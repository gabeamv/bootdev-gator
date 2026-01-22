package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gabeamv/bootdev-gator/gatorcommand"
	"github.com/gabeamv/bootdev-gator/internal/database"
	"github.com/gabeamv/bootdev-gator/internal/gatorconfig"
	_ "github.com/lib/pq"
)

func main() {
	config, err := gatorconfig.Read()
	if err != nil {
		fmt.Printf("error reading the configuration file: %v\n", err)
	}
	db, err := sql.Open("postgres", config.DBUrl) // Open a database connection
	if err != nil {
		fmt.Printf("error opening a database connection to %v: %v\n", config.DBUrl, err)
	}
	dbQueries := database.New(db)                          // we get all the queries that we have converted from sql to golang.
	state := gatorcommand.State{Db: dbQueries, S: &config} // state will be able to query the database as well as update the configuration file
	commands := gatorcommand.Commands{Commands: make(map[string]func(s *gatorcommand.State, c gatorcommand.Command) error)}

	commands.Register(gatorcommand.LOGIN, gatorcommand.HandlerLogin)
	commands.Register(gatorcommand.REGISTER, gatorcommand.HandlerRegister)
	commands.Register(gatorcommand.RESET, gatorcommand.HandlerReset)
	commands.Register(gatorcommand.USERS, gatorcommand.HandlerUsers)
	commands.Register(gatorcommand.AGG, gatorcommand.HandlerAgg)
	commands.Register(gatorcommand.ADDFEED, gatorcommand.MiddlewareLoggedIn(gatorcommand.HandlerAddFeed))
	commands.Register(gatorcommand.FEEDS, gatorcommand.HandlerFeeds)
	commands.Register(gatorcommand.FOLLOW, gatorcommand.MiddlewareLoggedIn(gatorcommand.HandlerFollow))
	commands.Register(gatorcommand.FOLLOWING, gatorcommand.MiddlewareLoggedIn(gatorcommand.HandlerFollowing))
	commands.Register(gatorcommand.UNFOLLOW, gatorcommand.MiddlewareLoggedIn(gatorcommand.HandlerUnfollow))
	commands.Register(gatorcommand.BROWSE, gatorcommand.MiddlewareLoggedIn(gatorcommand.HandlerBrowse))

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
