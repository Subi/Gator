package main

import (
	"context"
	"fmt"
	"github.com/subi/gator/internal/database"
	"strconv"
)

func handlerBrowse(s *state, cmd Command, user database.User) error {
	limit := 2
	if len(cmd.args) == 1 {
		i, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
		limit = i
	}
	posts, err := s.db.GetUserPost(context.Background(), database.GetUserPostParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}
	for _, post := range posts {
		fmt.Printf("%s\n%s\n%s\n", post.Title, post.Description, post.PublishedAt)
	}
	return nil
}
