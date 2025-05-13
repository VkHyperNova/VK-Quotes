package cmd

import (
	"fmt"
	"os"
	"strings"
	"vk-quotes/pkg/config"
	"vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func CommandLine(quotes *db.Quotes) {
	for {
		quotes.PrintCLI()

		var input string = ""
		var inputID int = 0

		fmt.Print("=> ")

		fmt.Scanln(&input, &inputID)

		input = strings.ToLower(input)

		switch input {

		case "add", "a":
			quotes.Add()
		case "update", "u":
			quotes.Update(inputID)
		case "delete", "d":
			quotes.Delete(inputID)
		case "find", "f":
			quotes.Find()
		case "showall", "s":
			quotes.PrintAllQuotes()
		case "stats":
			quotes.PrintStatistics()
		case "resetids":
			quotes.ResetIDs(quotes)
		case "read", "r":
			quotes.Read()
		case "similarquotes", "sim":
			db.FindSimilarQuotes(quotes)
		case "q", "quit":
			quotes.Backup()
			util.ClearScreen()
			os.Exit(0)
		default:
			config.AddMessage("Enter pressed!")
		}
	}
}
