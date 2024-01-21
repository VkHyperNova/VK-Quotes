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
			quote, author, language := Ask()
			Add(quote, author, language, DatabasePath)
			ReturnToCMD()
		case "update", "u":
			quote, author, language := Ask()
			Update(id, quote, author, language, DatabasePath)
			ReturnToCMD()
		case "delete", "d":
			Delete(id, DatabasePath)
			ReturnToCMD()
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

