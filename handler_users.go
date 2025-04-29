package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/subi/gator/internal/database"
	"log"
	"time"
)

func handlerRegister(s *state, cmd Command) error {
	if len(cmd.args) == 0 {
		log.Fatal("Please provide a username")
	}

	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	u, err := s.db.CreateUser(context.Background(), user)
	if err != nil {
		return err
	}
	err = s.config.SetUser(u.Name)
	if err != nil {
		return err
	}
	fmt.Printf("User %s has been created!\n", u.Name)
	fmt.Printf("%+v\n", u)
	return nil
}

func handlerLogin(s *state, cmd Command) error {
	if len(cmd.args) == 0 {
		log.Fatal("A username is required to login")
	}
	u, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	err = s.config.SetUser(u.Name)
	if err != nil {
		return err
	}
	fmt.Print(u)
	return nil
}

func handlerGetUsers(s *state, cmd Command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.config.CurrentUsername {
			fmt.Printf("* %s (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %s\n", user.Name)
	}
	return nil
}
