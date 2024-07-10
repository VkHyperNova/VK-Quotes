package main

import (
	"fmt"
	"os"
	"vk-quotes/pkg/cmd"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

var (
	ProgramVersion = "1.22"
	SaveFilePath   = "./database/quotes.json"
)

func main() {
	util.ClearScreen()
	util.CreateRequiredFiles(SaveFilePath)

	quotes := cmd.LoadQuotes(SaveFilePath)

	cmd.PrintCLI(&quotes, ProgramVersion)

	var command string = ""
	var id int = 0

	fmt.Scanln(&command, &id)

	for {
		switch command {
		case "add", "a":
			cmd.Add(&quotes, SaveFilePath)
			main()
		case "update", "u":
			cmd.Update(&quotes, id, SaveFilePath)
			main()
		case "delete", "d":
			cmd.Delete(&quotes, id, SaveFilePath)
			main()
		case "showall", "s":
			quotes.PrintQuotes()
			util.PressAnyKey()
			main()
		case "stats":
			cmd.PrintStatistics(&quotes)
			util.PressAnyKey()
			main()
		case "read", "r":
			quotes.FindIds()
			db.ReadMode = true
			cmd.PrintMessage = "<< Reading >>"
			main()
		case "q":
			util.ClearScreen()
			os.Exit(0)
		default:
			util.ClearScreen()
			if command != "" {
				db.ReadMode = false
				quotes.FindByAuthor(command)
				util.PressAnyKey()
			}

			/* Read Mode On */
			if db.ReadMode {
				if len(db.IDs) == 0 {
					cmd.PrintMessage = "<< You Have Read Everything! >>"
					db.ReadCounter = 0
					db.ReadMode = false
					quotes.GetLastId()
				}

				db.ReadCounter += 1
			}
			main()
		}
	}
}
