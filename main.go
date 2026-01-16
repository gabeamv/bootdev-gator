package main

import (
	"fmt"
	"os"

	"github.com/gabeamv/bootdev-gator/gatorcommand"
	"github.com/gabeamv/bootdev-gator/internal/gatorconfig"
)

func main() {
	config, err := gatorconfig.Read()
	if err != nil {
		fmt.Printf("error reading the configuration file: %v\n", err)
	}
	state := gatorcommand.State{S: &config}
	commands := gatorcommand.Commands{Commands: make(map[string]func(s *gatorcommand.State, c gatorcommand.Command) error)}

	commands.Register(gatorcommand.LOGIN, gatorcommand.HandlerLogin)

	commandName, args, err := CleanInput(os.Args)
	if err != nil {
		fmt.Printf("error cleaning arguments: %v\n", err)
		os.Exit(1)
	}
	command := gatorcommand.Command{Name: commandName, Args: args}
	err = commands.Run(&state, command)
	if err != nil {
		fmt.Printf("error: failed running command=%v args=%#v: %v\n", command.Name, command.Args, err)
		os.Exit(1)
	}

}

func CleanInput(args []string) (string, []string, error) {
	if len(args) < 2 {
		return "", []string{}, fmt.Errorf("error: args len=%v, arg=%#v, expected argument='command' + 'args'", len(args), args)
	}
	commandName := args[1]
	var commandArgs []string
	if len(args) > 2 {
		commandArgs = args[2:]
	}
	return commandName, commandArgs, nil
}

func TestConfig() {
	config, err := gatorconfig.Read()
	if err != nil {
		fmt.Printf("error reading the configuration file: %v\n", err)
	}
	err = config.SetUser("jamaal")
	if err != nil {
		fmt.Printf("error setting user: %v\n", err)
	}
	config, err = gatorconfig.Read()
	if err != nil {
		fmt.Printf("error reading the configuration file: %v\n", err)
	}
	fmt.Printf("Updated config file: %+v\n", config)
}
