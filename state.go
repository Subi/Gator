package main

import (
	"github.com/subi/gator/internal/config"
	"github.com/subi/gator/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}
