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
	// db.DATABASE = db.LoadDB(DatabasePath)
	Database := db.LoadDB(DatabasePath)

	PrintVKQUOTES(Version, LastItemIndex, &Database)


	var cmd string = ""
	var id int = 0

	fmt.Scanln(&cmd, &id)

	

	for {
		switch cmd {
		case "add", "a":
			quote, author, language := GetQuoteDetails(&Database)
			id := db.FindUniqueID(&Database)
			LastItemIndex = CheckIndex(id, &Database)
			Add(id, quote, author, language, &Database)
			db.SaveDB(&Database, DatabasePath)
			ReturnToCMD()
		case "update", "u":
			quote, author, language := GetQuoteDetails(&Database)
			LastItemIndex = CheckIndex(id, &Database)
			Update(LastItemIndex, quote, author, language, DatabasePath, &Database)
			db.SaveDB(&Database, DatabasePath)
			ReturnToCMD()
		case "delete", "d":
			LastItemIndex = CheckIndex(id, &Database)
			Delete(LastItemIndex, DatabasePath, &Database)
			db.SaveDB(&Database, DatabasePath)
			ReturnToCMD()
		case "showall", "s":
			PrintAllQuotes(&Database)
		case "stats":
			PrintStatistics(&Database)
		case "q":
			util.ClearScreen()
			os.Exit(0)
		default:
			util.ClearScreen()
			CMD()
		}
	}
}

func CheckIndex(id int, Database *[]db.Quotes) int {

	index := db.GetIndexFromId(id, Database)

	if index == -1 {
		util.PrintRed("\nIndex out of range!\n")
		ReturnToCMD()
	}

	// *LastItemIndex = index

	return index
}
