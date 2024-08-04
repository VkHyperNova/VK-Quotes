package main

import (
	"vk-quotes/pkg/cmd"
	"vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func main() {
	quotes := db.Quotes{}
	util.CreateDirectory()
	quotes.ReadFromFile()
	cmd.CommandLine(&quotes)
}

