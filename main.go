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
	if err := util.CreateFilesAndFolders(); err != nil {
		fmt.Println("Error creating files/folders:", err)
		os.Exit(1)
	}

	quotes := db.Quotes{}

	err := quotes.ReadFromFile(config.LocalFile)
	if err != nil {
		log.Fatalf("Fatal error: failed to load walkings database: %v", err)
	}

	cmd.CommandLine(&quotes)
}
