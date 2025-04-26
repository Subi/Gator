package main

import (
	"errors"
	"fmt"
)

type Command struct {
	name string
	args []string
}

type Commands struct {
	opts map[string]func(*state, Command) error
}

func handlerLogin(s *state, cmd Command) error {
	if len(cmd.args) == 0 {
		return errors.New("A username is required to login")
	}
	err := s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("Username has been set!")
	return nil
}

func (c *Commands) run(s *state, cmd Command) error {
	err := c.opts[cmd.name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) register(name string, f func(*state, Command) error) {
	c.opts[name] = f
}
