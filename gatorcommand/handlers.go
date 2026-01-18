package gatorcommand

import (
	"context"
	"fmt"
	"time"

	"github.com/gabeamv/bootdev-gator/gatorfeed"
	"github.com/gabeamv/bootdev-gator/internal/database"
	"github.com/google/uuid"
)

// Handler function to handle the login command.
func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("arguments expected: 'username'")
	}
	username := cmd.Args[0]
	user, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("error getting the user '%v': %w", username, err)
	}
	err = s.S.SetUser(user.Name) // arg 0 is string of username
	if err != nil {
		return fmt.Errorf("error setting the user as '%v': %w", user.Name, err)
	}
	fmt.Printf("The user has been set to %v.\n", user.Name)
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("arguments expected: 'username'")
	}
	username := cmd.Args[0]
	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), Name: username})
	if err != nil {
		return fmt.Errorf("error creating user=%v: %w", username, err)
	}
	s.S.SetUser(user.Name)
	fmt.Printf("User '%v' has been created.\n", user.Name)
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("error, no arguments expected")
	}
	err := s.Db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting all the registered users: %w", err)
	}
	fmt.Println("Deleted all registered users.")
	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("error, no arguments expected")
	}
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}
	var out string
	for _, user := range users {
		name := user.Name + "\n"
		if user.Name == s.S.CurrentUsername {
			name = user.Name + " (current)\n"
		}
		out += name
	}
	fmt.Print(out)
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	feed, err := gatorfeed.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error getting the feed: %w", err)
	}
	fmt.Println(feed)
	return nil
}
