package main

import "fmt"

type command struct {
	name      string
	arguments []string
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
