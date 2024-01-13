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
		case "add", "a":
			Add()
		case "update", "u":
			Update(id)
		case "delete", "d":
			Delete(id)
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
