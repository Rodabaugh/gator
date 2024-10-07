package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Rodabaugh/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}

	name := cmd.arguments[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user %s does not exist", name)
		} else {
			// Handle unexpected database errors.
			return err
		}
	}

	err = s.config.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Println("User switched successfully!")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}
	userName := cmd.arguments[0]

	_, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		if err == sql.ErrNoRows {
			// User does not exist, which is good for user creation.
		} else {
			// Handle unexpected database errors.
			return err
		}
	} else {
		return fmt.Errorf("user %s already exists", userName)
	}

	currentTime := time.Now()
	result, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(), CreatedAt: currentTime, UpdatedAt: currentTime, Name: userName,
	})
	if err != nil {
		return err
	}

	// Update current user
	err = handlerLogin(s, command{arguments: []string{result.Name}})
	if err != nil {
		return fmt.Errorf("failed to login user: %w", err)
	}

	fmt.Println("User created and logged in:", result.Name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("The database has been reset")
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}
