package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Rodabaugh/gator/internal/config"
	"github.com/Rodabaugh/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

func main() {
	// Initialize config
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// Connect to the database
	dbURL := cfg.DbURL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}
	dbQueries := database.New(db)
	appState := &state{
		db:     dbQueries,
		config: cfg,
	}

	if len(os.Args) < 2 {
		log.Fatalf("not enough arguments provided")
	}

	// Extract command name and arguments
	commandName := os.Args[1]
	commandArgs := os.Args[2:]

	// Create command instance
	cmd := command{
		name:      commandName,
		arguments: commandArgs,
	}

	// Get the commands map
	cmds := &commands{handlers: make(map[string]func(*state, command) error)}

	// Register commands
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerListFeedFollows))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))
	cmds.register("help", handlerHelp)

	// Look up and run the handler
	if err := cmds.run(appState, cmd); err != nil {
		log.Fatal(err)
	}
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("unable to get user: %w", err)
		}
		return handler(s, cmd, user)
	}
}

func handlerHelp(s *state, cmd command) error {
	fmt.Println("Available Commands")
	fmt.Println("* register")
	fmt.Println("* login")
	fmt.Println("* users")
	fmt.Println("* addfeed")
	fmt.Println("* feeds")
	fmt.Println("* follow")
	fmt.Println("* following")
	fmt.Println("* unfollow")
	fmt.Println("* agg")
	fmt.Println("* browse")
	fmt.Println("* help")
	return nil
}
