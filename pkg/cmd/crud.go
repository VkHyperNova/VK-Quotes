package cmd

import (
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func Add() {
	Quote := util.GetInput("Quote: ")

	if db.CheckDublicates(Quote) != -1 {
		util.PrintRed("\n<< This quote is in the database >>\n")
		PrintQuote(db.CheckDublicates(Quote))
		util.PressAnyKey()
		util.ClearScreen()
		CMD()
	}
	Author := util.GetInput("Auhtor: ")
	if Author == "" {
		Author = "Unknown"
	}
	Language := util.GetInput("Language: ")

	NewQuote := db.CompileQuote(Quote, Author, Language)
	db.DATABASE = append(db.DATABASE, NewQuote)
	db.SaveDB("Quote added!")

	util.PressAnyKey()
	util.ClearScreen()
	CMD()
}

func Update(id int) {
	index := db.SearchIndexByID(id)

	if index == -1 {
		util.PrintRed("\nIndex out of range!\n")
	} else {
		PrintQuote(index)
		UpdatedQuote := util.GetInput("\nUpdate Quote: ")
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
		db.SaveDB("Quote updated!")
	}
	util.PressAnyKey()
	util.ClearScreen()
	CMD()
}

func Delete(id int) {

	index := db.SearchIndexByID(id)

	if index == -1 {
		util.PrintRed("\nIndex out of range!\n")
	} else {
		PrintQuote(index)
		db.DATABASE = append(db.DATABASE[:index], db.DATABASE[index+1:]...)
		db.SaveDB("Quote deleted!")
	}
	util.PressAnyKey()
	util.ClearScreen()
	CMD()
}
