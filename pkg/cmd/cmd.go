package cmd

import (
	"fmt"
	"os"
	"vk-quotes/pkg/config"
	"vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func CommandLine(quotes *db.Quotes) {

	quotes.PrintCLI()

	var input string = ""
	var inputID int = 0

	fmt.Print("=> ")

	fmt.Scanln(&input, &inputID)

	for {
		switch input {

		case "add", "a":
			emptyQuote := db.Quote{}
			if quotes.UserInput(emptyQuote) {
				quotes.Add()
			}
			CommandLine(quotes)
		case "update", "u":
			UpdateQuote := quotes.FindQuoteByID(inputID)
			if quotes.UserInput(UpdateQuote) {
				quotes.Update(inputID)
			}
			CommandLine(quotes)

		case "delete", "d":
			quotes.Delete(inputID)
			CommandLine(quotes)

		case "find", "f":
			quotes.Find()
			CommandLine(quotes)
			
		case "showall", "s":
			quotes.PrintAllQuotes()
			CommandLine(quotes)

		case "stats":
			quotes.PrintStatistics()
			CommandLine(quotes)

		case "resetids":
			quotes.ResetIDs(quotes)
			CommandLine(quotes)

		case "read", "r":
			quotes.Read()
			CommandLine(quotes)

		case "similarquotes", "sim":
			db.FindSimilarQuotes(quotes)
			CommandLine(quotes)

		case "q", "quit":
			util.ClearScreen()
			os.Exit(0)

		default:
			config.AddMessage("Enter pressed!")
			CommandLine(quotes)
		}
	}
}
