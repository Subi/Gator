package main

import (
	"fmt"
	"github.com/subi/gator/internal/config"
	"log"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	err = cfg.SetUser("justin")

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	fmt.Printf("Reading config file %+v\n", cfg)
}
