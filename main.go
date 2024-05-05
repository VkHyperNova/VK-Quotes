package main

import (
	"fmt"
	"os"
	"vk-quotes/pkg/cmd"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func Abort(input string) {
	if input == "q" {
		cmd.ErrorMsg = "<< previous action aborted by user. >>"
		main()
	}
}

func GetQuote(Database *[]db.Quotes, id int) string {

	if id > 0 {
		index := db.FindIndex(id, Database)
		util.PrintGreen("\n\"" + (*Database)[index].QUOTE + "\"\n")
	} else {
		util.PrintGreen("\n\"" + "Unknown" + "\"\n")
	}

	util.PrintCyan("Quote: ")
	input := util.ScanUserInput()

	Abort(input)

	if db.FindDublicates(input, Database) != -1 {
		cmd.ErrorMsg = "<< there are dublicates in database. >>"
		cmd.CurrentQuoteIndex = db.FindDublicates(input, Database)
		main()
	}

	return input
}

func GetAuthor(Database *[]db.Quotes, id int) string {

	if id > 0 {
		index := db.FindIndex(id, Database)
		util.PrintGreen("\n\"" + (*Database)[index].AUTHOR + "\"\n")
	} else {
		util.PrintGreen("\n\"" + "Unknown" + "\"\n")
	}

	util.PrintCyan("Author: ")
	input := util.ScanUserInput()

	Abort(input)
	return input
}

func GetLanguage(Database *[]db.Quotes, id int) string {

	if id > 0 {
		index := db.FindIndex(id, Database)
		util.PrintGreen("\n\"" + (*Database)[index].LANGUAGE + "\"\n")
	} else {
		util.PrintGreen("\n\"" + "English" + "\"\n")
	}

	util.PrintCyan("Language: ")
	input := util.ScanUserInput()

	Abort(input)
	return input
}

func userInput(Database *[]db.Quotes, id int) []string {
	quote := GetQuote(Database, id)
	author := GetAuthor(Database, id)
	language := GetLanguage(Database, id)
	return []string{quote, author, language}
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
			input := userInput(&Database, 0)
			cmd.Create(input, &Database, cmd.DatabasePath)
			cmd.AddCount += 1
			cmd.ReadCount = 1
			main()
		case "update", "u":
			input := userInput(&Database, id)
			cmd.Update(id, input, &Database, cmd.DatabasePath)
			cmd.ReadCount = 1
			main()
		case "delete", "d":
			cmd.Delete(id, &Database, cmd.DatabasePath)
			cmd.ReadCount = 1
			main()
		case "showall", "s":
			cmd.PrintAllQuotes(&Database)
			cmd.ReadCount = 1
			util.PressAnyKey()
			main()
		case "stats":
			cmd.PrintStatistics(&Database)
			cmd.ReadCount = 1
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
			cmd.CurrentQuoteIndex = -1
			cmd.ReadCount += 1
			main()
		}
	}
}
