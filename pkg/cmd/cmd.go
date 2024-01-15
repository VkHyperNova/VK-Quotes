package cmd

import (
	"fmt"
	"os"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

var Version = "1.0"
var DatabasePath = "./database/quotes.json"

func CMD() {

	db.ValidateRequiredFiles(DatabasePath)
	db.DATABASE = db.LoadDB(DatabasePath)

	PrintVKQUOTES(Version)

	var cmd string = ""
	var id int = 0

	fmt.Scanln(&cmd, &id)

	for {
		switch cmd {
		case "add", "a":
			Quote, Author, Language := Ask()
			Add(Quote, Author, Language, DatabasePath)
			util.PressAnyKey()
			util.ClearScreen()
			CMD()
		case "update", "u":
			Update(id, DatabasePath)
		case "delete", "d":
			Delete(id, DatabasePath)
		case "showall", "s":
			PrintAllQuotes()
		case "stats":
			PrintStatistics()
		case "q":
			util.ClearScreen()
			os.Exit(0)
		default:
			util.ClearScreen()
			CMD()
		}
	}
}
