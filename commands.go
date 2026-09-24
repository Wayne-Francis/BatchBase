package main

import (
	"errors"

	"github.com/Wayne_Francis/BatchBase/internal/config"
	"github.com/Wayne_Francis/BatchBase/internal/database"
)

type command struct {
	Name string
	Args []string
}

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return f(s, cmd)
}
