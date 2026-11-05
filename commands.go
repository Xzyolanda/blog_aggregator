package main

import (
	"errors"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	AllCommands    map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.AllCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.AllCommands[name] = f
}

