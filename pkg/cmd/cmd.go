package cmd

import (
	"fmt"
	"os"
	"strconv"
	"vk-quotes/pkg/database"
	"vk-quotes/pkg/global"
	"vk-quotes/pkg/print"
)

func CMD() {

	global.DB = database.LoadQuotesDatabase()

	print.PrintQuotes()

	print.PrintProgramStart()
	print.PrintCommands()

	print.AddBrackets(strconv.Itoa(len(global.DB)))
	print.PrintCyan("=> ")

	var cmd string = ""
	var id int = 0

	fmt.Scanln(&cmd, &id)

	for {
		switch cmd {
		case "add":
			database.AddQuote()
			CMD()
		case "update":
			database.UpdateQuote(id)
			CMD()
		case "delete":
			database.DeleteQuote(id)
			CMD()
		case "q":
			print.ClearScreen()
			os.Exit(0)
		default:
			print.ClearScreen()
			CMD()
		}
	}
}
