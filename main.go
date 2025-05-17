package main

import (
	"vk-quotes/pkg/cmd"
	"vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

// main is the entry point of the application.
func main() {

	// It initializes a Quotes database
	quotes := db.Quotes{}

	// Ensures the necessary directory structure exists
	util.CreateDirectory()

	// Loads quotes from a file
	quotes.ReadFromFile()

	// and starts the command-line interface for user interaction.
	cmd.CommandLine(&quotes)
}
