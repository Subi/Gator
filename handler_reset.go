package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd Command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Print("Users data has been deleted")
	return nil
}
