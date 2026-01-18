package gatorcommand

import (
	"context"
	"fmt"
	"time"

	"github.com/gabeamv/bootdev-gator/gatorfeed"
	"github.com/gabeamv/bootdev-gator/internal/database"
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
	now := time.Now().UTC()
	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{CreatedAt: now, UpdatedAt: now, Name: username})
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

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("error, expecting args='name','url'")
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	currUser, err := s.Db.GetUser(context.Background(), s.S.CurrentUsername)
	if err != nil {
		return fmt.Errorf("error getting the current user: %w", err)
	}
	now := time.Now().UTC()
	feed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{CreatedAt: now, UpdatedAt: now, Name: name, Url: url, UserID: currUser.ID})
	if err != nil {
		return fmt.Errorf("error adding feed '%v' to the database: %w", name, err)
	}
	fmt.Printf("User '%v' has added feed: %v\n", currUser.Name, feed.Name)
	return nil
}
