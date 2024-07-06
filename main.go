package main

import (
	"fmt"
	"os"
	"time"
	"vk-quotes/pkg/cmd"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func main() {
	util.ClearScreen()
	util.CreateRequiredFiles(cmd.IsSaveFilePath)

	quotes := db.Quotes{}
	err := quotes.LoadQuotes("./database/quotes.json")
	if err != nil {
		fmt.Println("Error loading quotes:", err)
	}
	cmd.PrintCLI(&quotes)

	var command string = ""
	var id int = 0

	fmt.Scanln(&command, &id)

	for {
		switch command {
		case "add", "a":
			cmd.IsReadMode = false
			inputs, validation := cmd.UserInput(&quotes)
			if validation {
				newID := quotes.NewID()
				quotes.AddQuote(db.Quote{ID: newID, QUOTE: inputs[0], AUTHOR: inputs[1], LANGUAGE: inputs[2], DATE: time.Now().Format("02.01.2006")})
				quotes.SaveQuotes(cmd.IsSaveFilePath)
				cmd.IsMessage = fmt.Sprintf("<< %d Quote Added! >>", newID)
			}
			main()
		case "update", "u":
			cmd.IsReadMode = false
			updatedInputs := cmd.EditUserInput(&quotes, id)
			quotes.UpdateQuote(db.Quote{ID: id, QUOTE: updatedInputs[0], AUTHOR: updatedInputs[1], LANGUAGE: updatedInputs[2], DATE: time.Now().Format("02.01.2006")})
			quotes.SaveQuotes(cmd.IsSaveFilePath)
			cmd.IsMessage = fmt.Sprintf("<< %d Quote Updated! >>", id)
			main()
		case "delete", "d":
			cmd.IsReadMode = false
			quotes.DeleteQuote(id)
			quotes.SaveQuotes(cmd.IsSaveFilePath)
			cmd.IsMessage = fmt.Sprintf("<< %d Quote Deleted! >>", id)
			main()
		case "showall", "s":
			cmd.PrintAllQuotes(&quotes)
			util.PressAnyKey()
			main()
		case "stats":
			cmd.PrintStatistics(&quotes)
			util.PressAnyKey()
			main()
		case "read", "r":
			cmd.IsReadMode = true
			cmd.IsMessage = "<< Reading >>"
			cmd.MustPrintQuoteID = -1
			main()
		case "q":
			util.ClearScreen()
			os.Exit(0)
		default:
			util.ClearScreen()
			if command != "" {
				cmd.IsReadMode = false
				quotes.FindByAuthor(command)
				util.PressAnyKey()
			}
			
			/* Read Mode On */
			if cmd.IsReadMode {
				cmd.DeleteUsedIndexes(&quotes)
				cmd.IsReadCount += 1
			}
			main()
		}
	}
}
