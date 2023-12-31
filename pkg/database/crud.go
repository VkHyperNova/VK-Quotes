package database

import (
	"vk-quotes/pkg/global"
	"vk-quotes/pkg/print"
	"vk-quotes/pkg/util"
)

func AddQuote() {
	Quote := util.GetInput("Quote: ")
	Author := util.GetInput("Auhtor: ")
	if Author == "" {
		Author = "Unknown"
	}
	Language := util.GetInput("Language: ")

	NewQuote := util.CompileQuote(Quote, Author, Language)
	global.DB = append(global.DB, NewQuote)
	SaveQuoteDatabase()
}

func UpdateQuote(id int) {
	index := util.SearchIndexByID(id)
	confirm := false

	if index == -1 {
		print.PrintRed("\nIndex out of range!\n")
	} else {
		print.PrintQuote(index)
		confirm = util.Confirm()
	}

	if confirm {
		UpdatedQuote := util.GetInput("Update Quote: ")
		if UpdatedQuote == "" {
			UpdatedQuote = global.DB[index].QUOTE
		}
		UpdatedAuthor := util.GetInput("Update Author: ")
		if UpdatedAuthor == "" {
			UpdatedAuthor = global.DB[index].AUTHOR
		}
		UpdatedLanguage := util.GetInput("Update Language: ")
		if UpdatedLanguage == "" {
			UpdatedLanguage = global.DB[index].LANGUAGE
		}
		global.DB[index].QUOTE = UpdatedQuote
		global.DB[index].AUTHOR = UpdatedAuthor
		global.DB[index].LANGUAGE = UpdatedLanguage
		SaveQuoteDatabase()
		print.PrintGreen("Quote Updated!\n\n")
	} else {
		print.PrintGreen("Returning../\n\n")
	}

	print.ClearScreen()
}

func DeleteQuote(id int) {
	confirm := false

	index := util.SearchIndexByID(id)

	if index == -1 {
		print.PrintRed("\nIndex out of range!\n")
	} else {
		print.PrintQuote(index)
		confirm = util.Confirm()
	}

	if confirm {
		global.DB = append(global.DB[:index], global.DB[index+1:]...)
		SaveQuoteDatabase()
		print.PrintGreen("Quote deleted!\n\n")
	} else {
		print.PrintGreen("Returning../\n\n")
	}

	print.ClearScreen()
}
