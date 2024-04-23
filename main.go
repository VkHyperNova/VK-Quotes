package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"vk-quotes/pkg/cmd"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func userInput(Database *[]db.Quotes) []string {
	quoteFields := []string{"Quote: ", "Author: ", "Language: "}
	var inputs []string

	for _, field := range quoteFields {

		util.PrintCyan(field)

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		userInput := scanner.Text()
		userInput = strings.TrimSpace(userInput)

		if userInput == "q" {
			cmd.Msg = "<< previous action aborted by user. >>"
			main()
		}

		if field == "Quote: " {
			if db.FindDublicates(userInput, Database) != -1 {
				cmd.Msg = "<< there are dublicates in database. >>"
				cmd.CurrentQuoteIndex = db.FindDublicates(userInput, Database)
				main()
			}
		}

		inputs = append(inputs, userInput)
	}

	return inputs
}

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
			quote := userInput(&Database)
			cmd.Create(quote, &Database, cmd.DatabasePath)
			main()
		case "update", "u":
			quote := userInput(&Database)
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
			cmd.ReadCount += 1
			main()
		}
	}
}
