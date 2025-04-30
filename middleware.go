package main

import (
	"context"
	"github.com/subi/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd Command, user database.User) error) func(*state, Command) error {
	return func(s *state, cmd Command) error {
		u, err := s.db.GetUser(context.Background(), s.config.CurrentUsername)
		if err != nil {
			return err
		}
		return handler(s, cmd, u)
	}
}
