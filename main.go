package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/subi/gator/internal/config"
	"github.com/subi/gator/internal/database"
	"log"
	"os"
)

func main() {
	// Read environment variables from config file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	// Initiate db connection
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	dbQueries := database.New(db)
	defer db.Close()

	// Initialize state keep track of changes
	s := &state{
		db:     dbQueries,
		config: &cfg,
	}

	cmds := Commands{opts: make(map[string]func(*state, Command) error)}

	// Register handlers
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", initScrape)
	cmds.register("feeds", handlerListFeeds)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("follow", middlewareLoggedIn(handlerFollowFeed))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollowFeed))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	// Check if the length of arguments is at least 2
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
		log.Fatalf("Error running command: %v\n", err)
	}
}
