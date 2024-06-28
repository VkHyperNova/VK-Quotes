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
	cmd.PrintCLI(&Database)

	var command string = ""
	var id int = 0

	fmt.Scanln(&command, &id)

	for {
		switch command {
		case "add", "a":
			input, validation := cmd.UserInput(&Database, 0)
			if validation {
				cmd.Create(input, &Database, cmd.DatabasePath)
				cmd.AddCount += 1
			}
			cmd.ReadCount = 0
			main()
		case "update", "u":
			input, validation := cmd.UserInput(&Database, id)
			if validation {
				cmd.Update(id, input, &Database, cmd.DatabasePath)
			}
			cmd.ReadCount = 0
			main()
		case "delete", "d":
			cmd.Delete(id, &Database, cmd.DatabasePath)
			cmd.CurrentQuoteIndex = -1
			cmd.ReadCount = 0
			main()
		case "showall", "s":
			cmd.PrintAllQuotes(&Database)
			cmd.ReadCount = 0
			util.PressAnyKey()
			main()
		case "stats":
			cmd.PrintStatistics(&Database)
			cmd.ReadCount = 0
			util.PressAnyKey()
			main()
		case "q":
			util.ClearScreen()
			os.Exit(0)
		default:
			util.ClearScreen()
			if command != "" {
				cmd.FindByAuthor(&Database, command)
				util.PressAnyKey()
			}
			
			if cmd.CurrentQuoteIndex == -1 {
				cmd.ReadCount += 1
			}
			
			main()
		}
	}
}
