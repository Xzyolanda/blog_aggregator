package main

import (
	"context"
	"fmt"
)

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't return all users: %w", err)
	}
	for _, user := range users {
		if s.cfg.CurrentUserName == user.Name {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {fmt.Printf("* %v\n",user.Name)}
	}
	return nil
	
}