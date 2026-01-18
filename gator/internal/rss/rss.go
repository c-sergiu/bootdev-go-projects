package rss

import (
	"github.com/c-sergiu/bootdev-go-projects/gator/internal/database"
	"database/sql"
	"context"
	"encoding/xml"
	"net/http"
	"io"
	"html"
	"time"
	"fmt"
	"github.com/google/uuid"
	"strings"

)

import _ "github.com/lib/pq"

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err :=  http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("user-agent", "gator")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return  nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, _ := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}

func ScrapeFeeds(db *database.Queries) error {
	feedToFetch, err := db.GetNextFeedToFetch(context.Background())
	if err != nil {return err}

	err = db.MarkFeedFetched(
		context.Background(), 
		database.MarkFeedFetchedParams{
			ID: feedToFetch.ID,
			UpdatedAt: time.Now(),
		})
	if err != nil {return err}

	feed, err := FetchFeed(context.Background(), feedToFetch.Url)
	if err != nil {return err}

	for _, f := range feed.Channel.Item {
		pubTime, err := time.Parse(time.RFC1123Z, f.PubDate)
		var t sql.NullTime
		if err != nil {
			t = sql.NullTime{
				Valid: false,
			}
		} else {
			t = sql.NullTime{
				Valid: true,
				Time: pubTime,
			}
		}

		if f.Title == "" || f.Link == "" {
			continue
		}
		
		post, err := db.CreatePost(
			context.Background(),
			database.CreatePostParams{
				ID: uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				PublishedAt: t,
				Title: f.Title,
				Url: f.Link,
				Description: f.Description,
				FeedID: feedToFetch.ID,
			})

		if err != nil {
			errorToIgnore := "posts_url_key"
			ignore := strings.Contains(err.Error(), errorToIgnore) 
			if !ignore {
				return err
			} else {
				fmt.Printf("Ignored duplicate url %v\n", f.Link)
				continue
			}
		}

		fmt.Printf("Successfully inserted %v\n", post.Url)
	}
	return nil
}
