package gatorcommand

import (
	"errors"
	"fmt"

	"github.com/gabeamv/bootdev-gator/internal/gatorconfig"
)

type State struct {
	S *gatorconfig.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Commands map[string]func(s *State, cmd Command) error // map storing key value pairs of command name and its handler
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("arguments expected: 'username'")
	}
	username := cmd.Args[0]
	err := s.S.SetUser(username) // arg 0 is string of username
	if err != nil {
		return fmt.Errorf("error setting the user as '%v': %w", username, err)
	}
	fmt.Printf("The user has been set to %v.\n", username)
	return nil
}

func (c *Commands) Run(s *State, cmd Command) error {
	command, ok := c.Commands[cmd.Name]
	if !ok {
		return fmt.Errorf("error. no command exists of '%v'.", cmd.Name)
	}
	err := command(s, cmd)
	if err != nil {
		return fmt.Errorf("error running command '%v': %w", cmd.Name, err)
	}
	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) error {
	_, ok := c.Commands[name]
	if !ok {
		c.Commands[name] = f
		return nil
	}
	return fmt.Errorf("error. unable to register command '%v'.", name)
}
