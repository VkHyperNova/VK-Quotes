package main

import (
	"vk-quotes/pkg/cmd"
	"vk-quotes/pkg/config"
	"vk-quotes/pkg/db"
)

func main() {
	quotes := db.Quotes{}
	err := quotes.ReadFromFile(config.SaveQuotesPath, config.SaveFolder)
	if err != nil {	
		panic(err)
	}

	cmd.CMD(&quotes)
}
