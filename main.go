package main

import (
	"vk-quotes/pkg/cmd"
	"vk-quotes/pkg/db"
)

func main() {
	
	quotes := db.Quotes{}
	err := quotes.ReadFromFile()
	if err != nil {
		panic(err)
	}

	cmd.CommandLine(&quotes)
}
