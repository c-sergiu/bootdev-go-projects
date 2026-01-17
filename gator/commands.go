package main

import (
	"github.com/c-sergiu/bootdev-go-projects/gator/internal/database"
	"github.com/google/uuid"
	"time"
	"fmt"
	"context"
	"strconv"
)

import _ "github.com/lib/pq"

func handlerLogin(s* State, args[]string) error {
	if len(args) < 1 {
	return fmt.Errorf("No args")
	}

	user, err := s.db.GetUser(
	context.Background(),
		args[0])
	if err != nil {
		return err
	}

	if err := s.cfg.SetUser(user.Name); err != nil {
		return err
	}
	fmt.Println("Login Successful")
	return nil
}
func handlerRegister(s* State, args[]string) error {
	if len(args) < 1 {
		return fmt.Errorf("No Args")
	}

	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name: args[0],
		})

	if err != nil {
		return err
	}
	if err := s.cfg.SetUser(user.Name); err != nil {
		return err
	}
	fmt.Println("Register Successful")
	return nil
}
func handlerReset(s* State, args[]string) error {
	if err := s.db.DeleteAllUsers(context.Background()); err != nil {
		return err
	}
	return nil
}
func handlerUsers(s* State, args[]string) error {
	users, err := s.db.GetAllUsers(context.Background()); 
	if err != nil {
		return err
	}
	for _, user := range users {
		fmt.Printf("* %s", user.Name)
		if user.Name == s.cfg.CurUser {
			fmt.Printf(" (current)\n")
		} else {
			fmt.Println()
		}
	}
	return nil
}
func handlerAgg(s* State, args[]string) error {
	if len(args) < 1 {return fmt.Errorf("Not enaugh args")}
	timeBetween, err := time.ParseDuration(args[0])
	if err != nil {return err}

	ticker := time.NewTicker(timeBetween)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			return err
		}
	}

	return nil
}
func handlerAddFeed(s* State, args[]string, user database.User) error {
	if len(args) < 2 { return fmt.Errorf("Not enaugh args") }

	feedName := args[0]
	feedURL := args[1]
	
	feed_params := database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: feedName,
		Url: feedURL,
		UserID: user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), feed_params)
	if err != nil {
		return err
	}
	
	fmt.Printf("Feed created successfully!\nID: %v\nCreatedAt: %v\nUpdatedAt: %v\nName: %v\nUrl: %v\nUserID: %v\n",
		feed.ID,
		feed.CreatedAt,
		feed.UpdatedAt,
		feed.Name,
		feed.Url,
		feed.UserID)

	follow_params := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID: feed.ID,
		UserID: user.ID,
	}

	feed_follow, err := s.db.CreateFeedFollow(
		context.Background(),
		follow_params)
	if err != nil { return err }

	fmt.Printf("Successful follow!\nFeed: %s\nUser: %s\n",
		feed_follow.FeedName,
		feed_follow.UserName)

	return nil

}
func handlerFeeds(s* State, args[]string) error {
	feeds, err := s.db.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("Name: %s\nUrl: %s\nUser: %s\n",
			feed.Name,
			feed.Url,
			feed.UserName)
	}
	return nil

}
func handlerFollow(s* State, args[]string, user database.User) error {
	if len(args) < 1 {
		return fmt.Errorf("Not enaugh args")
	}
	url := args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil { return err }
	
	params := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID: feed.ID,
		UserID: user.ID,
	}

	feed_follow, err := s.db.CreateFeedFollow(
		context.Background(),
		params)
	if err != nil { return err }

	fmt.Printf("Successful follow!\nFeed: %s\nUser: %s\n",
		feed_follow.FeedName,
		feed_follow.UserName)

	return nil
}
func handlerFollowing(s* State, args[]string, user database.User) error {
	feed_follows, err := s.db.GetFeedFollowsForUser(
		context.Background(),
		user.ID)
	if err != nil {return err}
	for _, feed_follow := range feed_follows {
		fmt.Printf("%s\n", feed_follow.FeedName)
	}
	return nil
}

func handlerUnfollow(s* State, args[]string, user database.User) error {
	if len(args) < 1 {return fmt.Errorf("Not enaugh args")}
	url := args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {return err}

	if err := s.db.DeleteFeedFollow(
		context.Background(),
		database.DeleteFeedFollowParams{
			UserID: user.ID,
			FeedID: feed.ID,
		}); err != nil {return err}
	
	fmt.Println("Successful unfollow")

	return nil
}

func handlerBrowse(s *State, args[] string, user database.User) error {
	var limit int32
	if len(args) < 1 {
		limit = 2
	} else { 
		r, err := strconv.Atoi(args[0])
		if err != nil {return err}
		limit = int32(r)
	}
	
	posts, err := s.db.GetPostsForUser(
		context.Background(),
		database.GetPostsForUserParams{
			UserID: user.ID,
			Limit: limit,
		})
	if err != nil {
		return err
	}

	for _, p := range posts {
		fmt.Printf("Title: %s\nDescription: %s\nUrl: %s\n",
			p.Title,
			p.Description,
			p.Url)
	}

	return nil
}
