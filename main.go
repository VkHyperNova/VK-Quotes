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
	util.ValidateRequiredFiles(cmd.IsDatabasePath)
	Database := db.OpenDB(cmd.IsDatabasePath)
	cmd.PrintCLI(&Database)

	var command string = ""
	var id int = 0

	fmt.Scanln(&command, &id)

	for {
		switch command {
		case "add", "a":
			validation := cmd.UserInput(&Database, 0)
			if validation {
				cmd.Add(&Database, cmd.IsDatabasePath)
			}
			main()
		case "update", "u":
			cmd.EditUserInput(&Database, id)
			cmd.Update(id, &Database, cmd.IsDatabasePath)
			main()
		case "delete", "d":
			cmd.Delete(id, &Database, cmd.IsDatabasePath)
			main()
		case "showall", "s":
			cmd.PrintAllQuotes(&Database)
			util.PressAnyKey()
			main()
		case "stats":
			cmd.PrintStatistics(&Database)
			util.PressAnyKey()
			main()
		case "read", "r":
			cmd.Read()
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
			
			/* Read Mode On */
			if cmd.IsReadMode {
				cmd.DeleteUsedIndexes(&Database)
				cmd.IsReadCount += 1
			}

			main()
		}
	}
}