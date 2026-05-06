package main

import (
	"fmt"
	"log"
	"os"
	"vk-quotes/pkg/cmd"
	"vk-quotes/pkg/config"
	"vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func main() {
	if err := util.InitLocalStorage(); err != nil {
		fmt.Println("Error creating files/folders:", err)
		os.Exit(1)
	}

	q := db.Quotes{}

	err := q.LoadFromFile(config.LocalFile)
	if err != nil {
		log.Fatalf("Fatal error: failed to load walkings database: %v", err)
	}

	cmd.CommandLine(&q)
}
