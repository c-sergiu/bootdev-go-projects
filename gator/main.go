package main

import (
	"github.com/c-sergiu/bootdev-go-projects/gator/internal/config"
	"github.com/c-sergiu/bootdev-go-projects/gator/internal/repl"
	"github.com/c-sergiu/bootdev-go-projects/gator/internal/database"
	"log"
	"os"
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

type State struct {
	cfg *config.Config
	db *database.Queries
}

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

func main() {
	// Init cfg
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	// Init db
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	
	// Init State
	state := &State{
		cfg: cfg,
db: dbQueries,
	}
	
	// Init cmd
	repl := repl.NewREPL(state, "gator")
	registerCommands(repl)
	
	// Handle
	args := os.Args[1:]
	if err := repl.HandleCommand(args); err != nil {
		log.Fatal(err)
	}
}

func registerCommands(r *repl.REPL[*State]) {
	r.Register(repl.Command[*State]{
		Name: "login",
		Desc: "Login user",
		Exec: handlerLogin,
	})
	r.Register(repl.Command[*State]{
		Name: "register",
		Desc: "Register new user",
		Exec: handlerRegister,
	})
	r.Register(repl.Command[*State]{
		Name: "reset",
		Desc: "Resets all data",
		Exec: handlerReset,
	})
	r.Register(repl.Command[*State]{
		Name: "users",
		Desc: "List all users",
		Exec: handlerUsers,
	})
	r.Register(repl.Command[*State]{
		Name: "agg",
		Desc: "Run feed loop",
		Exec: handlerAgg,
	})
	r.Register(repl.Command[*State]{
		Name: "addfeed",
		Desc: "Add a feed and follow",
		Exec: middlewareLoggedIn(handlerAddFeed),
	})
	r.Register(repl.Command[*State]{
		Name: "feeds",
		Desc: "List feeds",
		Exec: handlerFeeds,
	})
	r.Register(repl.Command[*State]{
		Name: "follow",
		Desc: "Follow a feed",
		Exec: middlewareLoggedIn(handlerFollow),
	})
	r.Register(repl.Command[*State]{
		Name: "following",
		Desc: "List feeds followed",
		Exec: middlewareLoggedIn(handlerFollowing),
	})
	r.Register(repl.Command[*State]{
		Name: "unfollow",
		Desc: "Unfollow feed",
		Exec: middlewareLoggedIn(handlerUnfollow),
	})
	r.Register(repl.Command[*State]{
		Name: "browse",
		Desc: "Browse Posts",
		Exec: middlewareLoggedIn(handlerBrowse),
	})

}

func middlewareLoggedIn(handler func(s *State, args[]string, user database.User) error) func(*State, []string) error {
	return func(s *State, args []string) error {
		userName := s.cfg.CurUser
		user, err := s.db.GetUser(
			context.Background(),
			userName)
		if err != nil {
			return err
		}

		err = handler(s, args, user)

		return err
	}
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
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

func scrapeFeeds(s *State) error {
	feedToFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {return err}

	err = s.db.MarkFeedFetched(
		context.Background(), 
		database.MarkFeedFetchedParams{
			ID: feedToFetch.ID,
			UpdatedAt: time.Now(),
		})
	if err != nil {return err}

	feed, err := fetchFeed(context.Background(), feedToFetch.Url)
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
		
		post, err := s.db.CreatePost(
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
