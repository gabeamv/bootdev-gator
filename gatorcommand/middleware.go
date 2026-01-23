package gatorcommand

import (
	"context"
	"fmt"

	"github.com/gabeamv/bootdev-gator/internal/database"
)

func MiddlewareLoggedIn(handler func(*State, Command, database.User) error) func(*State, Command) error {
	return func(s *State, c Command) error {
		user, err := s.Db.GetUser(context.Background(), s.S.CurrentUsername)
		if err != nil {
			return fmt.Errorf("error getting the user '%v' from middlewareloggedin: %w", s.S.CurrentUsername, err)
		}
		err = handler(s, c, user)
		if err != nil {
			return fmt.Errorf("error running handler from middlewareloggedin: %w", err)
		}
		return nil
	}
}

func MiddlewareAllCommands(handler func(*State, Command, Commands) error, commands Commands) func(*State, Command) error {
	return func(s *State, c Command) error {
		err := handler(s, c, commands)
		if err != nil {
			return fmt.Errorf("error middleware commands: %w", err)
		}
		return nil
	}
}
