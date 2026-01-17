package gatorcommand

import (
	"fmt"

	"github.com/gabeamv/bootdev-gator/internal/database"
	"github.com/gabeamv/bootdev-gator/internal/gatorconfig"
)

type State struct {
	Db *database.Queries
	S  *gatorconfig.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Commands map[string]func(s *State, cmd Command) error // map storing key value pairs of command name and its handler
}

// The Commands struct runs a specific command with its own handler.
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

// Register a command with its name and its own handler.
func (c *Commands) Register(name string, f func(*State, Command) error) error {
	_, ok := c.Commands[name]
	if !ok {
		c.Commands[name] = f
		return nil
	}
	return fmt.Errorf("error. unable to register command '%v'.", name)
}
