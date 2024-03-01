package cmd

import (
	"fmt"
	"os"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

var Version = "1.0"
var DatabasePath = "./database/quotes.json"
var LastItemIndex = -1

func CMD() {

	db.ValidateRequiredFiles(DatabasePath)
	Database := db.ReadDB(DatabasePath)

	PrintVKQUOTES(Version, LastItemIndex, &Database)

	var cmd string = ""
	var id int = 0

	fmt.Scanln(&cmd, &id)

	for {
		switch cmd {
		case "add", "a":
			quote, author, language := GetQuoteDetails(&Database)
			newID := db.FindUniqueID(&Database)
			Add(newID, quote, author, language, &Database)
			db.SaveDB(&Database, DatabasePath)
			LastItemIndex = CheckIndex(newID, &Database)
			util.PressAnyKey()
			CMD()
		case "update", "u":
			quote, author, language := GetQuoteDetails(&Database)
			LastItemIndex = CheckIndex(id, &Database)
			Update(LastItemIndex, quote, author, language, DatabasePath, &Database)
			db.SaveDB(&Database, DatabasePath)
			util.PressAnyKey()
			CMD()
		case "delete", "d":
			LastItemIndex = CheckIndex(id, &Database)
			Delete(LastItemIndex, DatabasePath, &Database)
			db.SaveDB(&Database, DatabasePath)
			util.PressAnyKey()
			CMD()
		case "showall", "s":
			PrintAllQuotes(&Database)
			util.PressAnyKey()
			CMD()
		case "stats":
			PrintStatistics(&Database)
			util.PressAnyKey()
			CMD()
		case "q":
			util.ClearScreen()
			os.Exit(0)
		default:
			util.ClearScreen()
			if cmd != "" {
				Search(&Database, cmd)
				util.PressAnyKey()
			}
			CMD()
		}
	}
}

func CheckIndex(id int, Database *[]db.Quotes) int {

	index := db.GetIndexFromId(id, Database)

	if index == -1 {
		util.PrintRed("\nIndex out of range! (CheckIndex function)\n")
		util.PressAnyKey()
		CMD()
	}

	return index
}
