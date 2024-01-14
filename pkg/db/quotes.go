package db

import (
	"encoding/json"
	"strings"
	"vk-quotes/pkg/util"
)

var DATABASE []Quotes

var DatabasePath = "./database/quotes.json"
var LastAddID = -1

type Quotes struct {
	ID       int    `json:"id"`
	QUOTE    string `json:"quote"`
	AUTHOR   string `json:"author"`
	LANGUAGE string `json:"language"`
	DATE     string `json:"date"`
}

func SaveDB(action string) {
	DatabaseAsByte := util.InterfaceToByte(DATABASE)
	util.WriteDataToFile(DatabasePath, DatabaseAsByte)
	util.PrintGreen("\n" + action + "\n")
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
	uniqueID := FindUniqueID()
	LastAddID = uniqueID

	return Quotes{
		ID:       uniqueID,
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

func CheckDublicates(quote string) int {
	for _, value := range DATABASE {
		if strings.EqualFold(value.QUOTE, quote) {
			return SearchIndexByID(value.ID)
		}
	}
	return -1
}

func SortNames(s string) []string {

	var names []string

	for _, value := range DATABASE {

		field := value.LANGUAGE
		
		if s == "authors" {
			field = value.AUTHOR
		}

		if !util.Contains(names, field) {
			names = append(names, field)
		}

	}

	return names
}

func CountNames(s string, names []string) map[string]int {

	myMap := make(map[string]int)

	for _, name := range names {

		for _, value := range DATABASE {

			field := value.LANGUAGE

			if s == "authors" {
				field = value.AUTHOR
			}

			if field == name {
				myMap[name] += 1
			}	
		}
	}

	return myMap
}
