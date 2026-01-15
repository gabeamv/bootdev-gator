package main

import (
	"fmt"

	"github.com/gabeamv/bootdev-gator/internal/configure"
)

func main() {
	config, err := configure.Read()
	if err != nil {
		fmt.Printf("error reading the configuration file: %v\n", err)
	}
	err = config.SetUser("jamaal")
	if err != nil {
		fmt.Printf("error setting user: %v\n", err)
	}
	config, err = configure.Read()
	if err != nil {
		fmt.Printf("error reading the configuration file: %v\n", err)
	}
	fmt.Printf("Updated config file: %+v\n", config)
}
