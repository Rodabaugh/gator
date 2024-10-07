package main

import (
	"database/sql"
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
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", handlerFollow)
	cmds.register("following", handlerListFeedFollows)

	// Look up and run the handler
	if err := cmds.run(appState, cmd); err != nil {
		log.Fatal(err)
	}
}
