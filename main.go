package main

import (
	"fmt"
	"github.com/subi/gator/internal/config"
	"log"
	"os"
)

func main() {

	// Read environment variables from config file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	// Initialize state keep track of changes
	s := &state{config: &cfg}

	cmds := Commands{opts: make(map[string]func(*state, Command) error)}
	// Register handlers
	cmds.register("login", handlerLogin)

	// Check if length of arguments is at least 2
	if len(os.Args) < 2 {
		log.Fatal("Please provide a command")
	}

	cmd := Command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	// Process command values submitted by user
	err = cmds.run(s, cmd)
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
	}
}
