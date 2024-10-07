package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Rodabaugh/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.arguments) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.name)
	}
	name := cmd.arguments[0]
	url := cmd.arguments[1]

	currentUser, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return err
	}

	currentTime := time.Now()

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      name,
		Url:       url,
		UserID:    currentUser.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	err = handlerFollow(s, command{arguments: []string{feed.Url}})
	if err != nil {
		return err
	}

	fmt.Println("Successfully created feed:")
	printFeed(feed, currentUser)
	fmt.Println()
	fmt.Println("=======================================================")
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		feedUser, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("database error: was unable to find user\nID: %s", feed.UserID)
		}
		printFeed(feed, feedUser)
		fmt.Println("=======================================================")
	}
	return nil
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.name)
	}
	url := cmd.arguments[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	currentUser, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return err
	}

	currentTime := time.Now()

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't follow feed: %w", err)
	}

	fmt.Println("Successfully followed feed:")
	fmt.Printf("%s is now being followed by %s\n", feedFollow.FeedName, feedFollow.UserName)
	fmt.Println("=======================================================")
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	currentUser, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return err
	}

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("* %s\n", feed.Name)
	}
	fmt.Println("=======================================================")
	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
}
