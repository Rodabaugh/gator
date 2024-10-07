package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Rodabaugh/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, currentUser database.User) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.name)
	}
	url := cmd.arguments[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
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

func handlerListFeedFollows(s *state, cmd command, currentUser database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return err
	}

	if len(feeds) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}

	fmt.Println("Feeds following:")
	for _, feed := range feeds {
		printFeedFollow(currentUser.Name, feed.Name)
	}
	fmt.Println("=======================================================")
	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}
