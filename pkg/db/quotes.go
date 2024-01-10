package db

import (
	"encoding/json"
	"vk-quotes/pkg/util"
)

var DATABASE []Quotes

var DatabasePath = "./database/quotes.json"

type Quotes struct {
	ID       int    `json:"id"`
	QUOTE    string `json:"quote"`
	AUTHOR   string `json:"author"`
	LANGUAGE string `json:"language"`
	DATE     string `json:"date"`
}

func SaveDB() {
	DatabaseAsByte := util.InterfaceToByte(DATABASE)
	util.WriteDataToFile(DatabasePath, DatabaseAsByte)
}

func LoadDB() []Quotes {
	file := util.ReadFile(DatabasePath)
	data := GetQuotesArray(file)

	return data
}

func ValidateRequiredFiles() {
	if !util.DoesDirectoryExist(DatabasePath) {
		util.CreateDirectory("database")
		util.WriteDataToFile(DatabasePath, []byte("[]"))
		util.PrintRed("New Database Created!\n")
	}
}

func FindUniqueID() int {

	if len(DATABASE) == 0 {
		return 1
	}

	return DATABASE[len(DATABASE)-1].ID + 1
}

func CompileQuote(quote string, author string, language string) Quotes {

	return Quotes{
		ID:       FindUniqueID(),
		QUOTE:    quote,
		AUTHOR:   author,
		LANGUAGE: language,
		DATE:     util.GetFormattedDate(),
	}
}

func GetQuotesArray(body []byte) []Quotes {

	QuotesStruct := []Quotes{}

	err := json.Unmarshal(body, &QuotesStruct)
	util.HandleError(err)

	return QuotesStruct
}

func SearchIndexByID(id int) int {

	index := -1

	for key, website := range DATABASE {
		if id == website.ID {
			index = key
		}
	}

	return index
}