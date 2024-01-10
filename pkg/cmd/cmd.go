package cmd

import (
	"fmt"
	"os"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

var Version = "1.0"

func CMD() {

	db.DATABASE = db.LoadDB()

	PrintVKQUOTES(Version)

	var cmd string = ""
	var id int = 0

	fmt.Scanln(&cmd, &id)

	for {
		switch cmd {
		case "add":
			AddQuote()
			CMD()
		case "update":
			UpdateQuote(id)
			CMD()
		case "delete":
			DeleteQuote(id)
			CMD()
		case "showall":
			PrintAllQuotes()
			CMD()
		case "q":
			util.ClearScreen()
			os.Exit(0)
		default:
			util.ClearScreen()
			CMD()
		}
	}
}
