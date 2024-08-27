package cmd

import (
	"os"
	"vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func CommandLine(quotes *db.Quotes) {

	quotes.PrintCLI()

	input, inputID := util.CommandPrompt("> ")

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
			if input != "" {
				quotes.PrintQuote(input)
				util.PressAnyKey()
			}
			CommandLine(quotes)
		}
	}
}