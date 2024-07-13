package main

import (
	"fmt"
	"os"
	"vk-quotes/pkg/cmd"
	"vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

var (
	ProgramVersion = "1.23"
	SaveFilePath   = "./database/quotes.json"
)

func main() {
	util.ClearScreen()
	util.CreateRequiredFiles(SaveFilePath)

	quotes := db.Quotes{}
	err := quotes.ReadFromFile(SaveFilePath)
	if err != nil {
		fmt.Println("Error loading quotes:", err)
	}

	cmd.PrintCLI(&quotes, ProgramVersion)

	var command string = ""
	var id int = 0

	fmt.Scanln(&command, &id)

	for {
		switch command {
		case "add", "a":
			validation := quotes.UserInput(0)
			if validation {
				cmd.Add(&quotes, util.UserInputs, SaveFilePath)
			}
			main()
		case "update", "u":
			validation := quotes.UserInput(id)
			if validation {
				cmd.Update(&quotes, util.UserInputs, id, SaveFilePath)
			}
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
			util.ReadMode = true
			util.Message = "<< Reading >>"
			main()
		case "q":
			util.ClearScreen()
			os.Exit(0)
		default:
			util.ClearScreen()
			if command != "" {
				util.ReadMode = false
				quotes.FindByAuthor(command)
				util.PressAnyKey()
			}

			/* Read Mode On */
			if util.ReadMode {
				if len(util.IDs) == 0 {
					util.Message = "<< You Have Read Everything! >>"
					util.ReadCounter = 0
					util.ReadMode = false
					quotes.FindLastId()
				}

				util.ReadCounter += 1
			}
			main()
		}
	}
}
