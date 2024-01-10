package cmd

import (
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func AddQuote() {
	Quote := util.GetInput("Quote: ")
	Author := util.GetInput("Auhtor: ")
	if Author == "" {
		Author = "Unknown"
	}
	Language := util.GetInput("Language: ")

	NewQuote := db.CompileQuote(Quote, Author, Language)
	db.DATABASE = append(db.DATABASE, NewQuote)
	db.SaveDB()
}

func UpdateQuote(id int) {
	index := db.SearchIndexByID(id)
	confirm := false

	if index == -1 {
		util.PrintRed("\nIndex out of range!\n")
	} else {
		PrintQuote(index)
		confirm = util.Confirm()
	}

	if confirm {
		UpdatedQuote := util.GetInput("Update Quote: ")
		if UpdatedQuote == "" {
			UpdatedQuote = db.DATABASE[index].QUOTE
		}
		UpdatedAuthor := util.GetInput("Update Author: ")
		if UpdatedAuthor == "" {
			UpdatedAuthor = db.DATABASE[index].AUTHOR
		}
		UpdatedLanguage := util.GetInput("Update Language: ")
		if UpdatedLanguage == "" {
			UpdatedLanguage = db.DATABASE[index].LANGUAGE
		}
		db.DATABASE[index].QUOTE = UpdatedQuote
		db.DATABASE[index].AUTHOR = UpdatedAuthor
		db.DATABASE[index].LANGUAGE = UpdatedLanguage
		db.SaveDB()
		util.PrintGreen("Quote Updated!\n\n")
	} else {
		util.PrintGreen("Returning../\n\n")
	}

	util.ClearScreen()
}

func DeleteQuote(id int) {
	confirm := false

	index := db.SearchIndexByID(id)

	if index == -1 {
		util.PrintRed("\nIndex out of range!\n")
	} else {
		PrintQuote(index)
		confirm = util.Confirm()
	}

	if confirm {
		db.DATABASE = append(db.DATABASE[:index], db.DATABASE[index+1:]...)
		db.SaveDB()
		util.PrintGreen("Quote deleted!\n\n")
	} else {
		util.PrintGreen("Returning../\n\n")
	}

	util.ClearScreen()
}
