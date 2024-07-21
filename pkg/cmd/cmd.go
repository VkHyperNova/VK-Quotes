package cmd

import (
	"fmt"
	"os"
	"time"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func CMD(quotes *db.Quotes, settings *util.Settings) {

	PrintCLI(quotes, settings)

	settings.Command, settings.ID = util.CommandPrompt(settings)

	for {
		switch settings.Command {
		case "add", "a":
			validation := quotes.UserInput(settings)
			if validation {
				Add(quotes, settings)
			}
			CMD(quotes, settings)
		case "update", "u":
			validation := quotes.UserInput(settings)
			if validation {
				Update(quotes, settings)
			}
			CMD(quotes, settings)
		case "delete", "d":
			Delete(quotes, settings)
			CMD(quotes, settings)
		case "showall", "s":
			quotes.PrintQuotes()
			util.PressAnyKey()
			CMD(quotes, settings)
		case "stats":
			PrintStatistics(quotes)
			util.PressAnyKey()
			CMD(quotes, settings)
		case "read", "r":
			Read(quotes, settings)
			CMD(quotes, settings)
		case "similar":
			db.RunTaskWithProgress(quotes, settings)
			util.PressAnyKey()
			CMD(quotes, settings)
		case "q":
			util.ClearScreen()
			os.Exit(0)
		default:
			if settings.Command != "" {
				quotes.PrintByAuthor(settings.Command)
				util.PressAnyKey()
			}
			CMD(quotes, settings)
		}
	}

}

func Add(quotes *db.Quotes, settings *util.Settings) bool {
	newID := quotes.CreateId()
	quotes.Add(db.Quote{ID: newID, QUOTE: settings.UserInputs[0], AUTHOR: settings.UserInputs[1], LANGUAGE: settings.UserInputs[2], DATE: time.Now().Format("02.01.2006")})
	quotes.SaveToFile(settings)
	settings.Message = fmt.Sprintf("<< %d Quote Added! >>", newID)

	return true
}

func Update(quotes *db.Quotes, settings *util.Settings) bool {

	quotes.Update(db.Quote{ID: settings.ID, QUOTE: settings.UserInputs[0], AUTHOR: settings.UserInputs[1], LANGUAGE: settings.UserInputs[2], DATE: time.Now().Format("02.01.2006")})
	quotes.SaveToFile(settings)
	settings.Message = fmt.Sprintf("<< %d Quote Updated! >>", settings.ID)

	return true
}

func Delete(quotes *db.Quotes, settings *util.Settings) bool {
	quotes.Delete(settings)
	quotes.SaveToFile(settings)
	settings.Message = fmt.Sprintf("<< %d Quote Deleted! >>", settings.ID)
	return true
}

func Read(quotes *db.Quotes, settings *util.Settings) {

	quotes.AppendRandomIDs(settings)
	settings.Message = "<< Reading Mode >>"

	for len(settings.RandomIDs) != 0 {
		PrintCLI(quotes, settings)
		settings.ReadCounter += 1
		var quit string
		fmt.Scanln(&quit)
		if quit == "q" {
			settings.Message = "<< Reading Mode Off >>"
			settings.ReadCounter = 0
			quotes.ResetID(settings)
			if len(settings.RandomIDs) > 0 {
				settings.RandomIDs = settings.RandomIDs[:0]
			}
			CMD(quotes, settings)
		}
	}
}
