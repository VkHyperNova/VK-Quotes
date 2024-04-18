package main

import (
	"fmt"
	"os"
	"vk-quotes/pkg/cmd"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func main() {
	util.ClearScreen()
	util.ValidateRequiredFiles(cmd.DatabasePath)
	Database := db.OpenDB(cmd.DatabasePath)

	cmd.PrintCLI(cmd.Version, cmd.CurrentQuoteIndex, &Database)

	var command string = ""
	var id int = 0

	fmt.Scanln(&command, &id)

	for {
		switch command {
		case "add", "a":
			quote := db.GetUserInput()
			cmd.Create(quote, &Database, cmd.DatabasePath)
			main()
		case "update", "u":
			quote := db.GetUserInput()
			db.ProcessUserInput(quote, &Database)
			cmd.Update(id, quote, &Database, cmd.DatabasePath)
			main()
		case "delete", "d":
			cmd.Delete(id, &Database, cmd.DatabasePath)
			main()
		case "showall", "s":
			cmd.PrintAllQuotes(&Database)
			util.PressAnyKey()
			main()
		case "stats":
			cmd.PrintStatistics(&Database)
			util.PressAnyKey()
			main()
		case "q":
			util.ClearScreen()
			os.Exit(0)
		default:
			util.ClearScreen()
			if command != "" {
				cmd.Read(&Database, command)
				util.PressAnyKey()
			}
			main()
		}
	}
}
