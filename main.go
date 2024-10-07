package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Rodabaugh/gator/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	name      string
	arguments []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("must have a name")
	}
	s.config.SetUser(cmd.arguments[0])
	println("User has been set.")
	return nil
}

type commands struct {
	handlers map[string]func(*state, command) error
}

// Register a new command by name with a handler function
func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

// Run a command if it exists in the handlers map
func (c *commands) run(s *state, cmd command) error {
	if handler, exists := c.handlers[cmd.name]; exists {
		return handler(s, cmd)
	}
	return fmt.Errorf("unknown command: %s", cmd.name)
}

func main() {
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

	// Initialize config and state
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	appState := &state{config: cfg}

	// Get the commands map
	cmds := &commands{handlers: make(map[string]func(*state, command) error)}

	// Register a command
	cmds.register("login", handlerLogin)

	// Look up and run the handler
	if err := cmds.run(appState, cmd); err != nil {
		log.Fatal(err)
	}
}
