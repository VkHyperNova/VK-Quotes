package db

import (
	"encoding/json"
	"strings"
	"vk-quotes/pkg/util"
)

type Quotes struct {
	ID       int    `json:"id"`
	QUOTE    string `json:"quote"`
	AUTHOR   string `json:"author"`
	LANGUAGE string `json:"language"`
	DATE     string `json:"date"`
}

func SaveDB(DB *[]Quotes, DatabasePath string) bool {
	util.WriteDataToFile(DatabasePath, util.StructToJson(DB))
	return true
}

func OpenDB(DatabasePath string) []Quotes {
	file := util.ReadFile(DatabasePath)

	Quotes := []Quotes{}

	err := json.Unmarshal(file, &Quotes)
	util.HandleError(err)

	return Quotes
}

func FindID(Database *[]Quotes) int {
	if len(*Database) == 0 {
		return 1
	}
	return (*Database)[len(*Database)-1].ID + 1
}

func GetAllNamesOf(s string, Database *[]Quotes) []string {

	var names []string

	for _, value := range *Database {

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
/* Print all the authors at once */
func Authors(Database *[]Quotes) {
	var authors []string
	var langauges []string

	for _, quote := range *Database {
		if !util.Contains(authors, quote.AUTHOR) {
			authors = append(authors, quote.AUTHOR)
		}

		if !util.Contains(langauges, quote.LANGUAGE) {
			langauges = append(langauges, quote.LANGUAGE)
		}
	}

	authorsByCount := make(map[string]int)
	
	for _, author := range authors {
		for _, quote := range *Database {
			if author == quote.AUTHOR {
				authorsByCount[quote.AUTHOR] += 1
			}
		}	
	}

	languagesByCount := make(map[string]int)

	for _, langauge := range langauges {
		for _, quote := range *Database {
			if langauge == quote.LANGUAGE {
				languagesByCount[quote.LANGUAGE] += 1
			}
		}	
	}

	

}

func CountByName(s string, names []string, Database *[]Quotes) map[string]int {

	myMap := make(map[string]int)

	for _, name := range names {
		for _, value := range *Database {
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

func FindDublicates(quote string, Database *[]Quotes) bool {

	for _, value := range *Database {
		if strings.EqualFold(value.QUOTE, quote) {
			return true
		}
	}

	return false
}

func GetUserInput() []string {
	quoteDetails := []string{"Quote: ", "Author: ", "Language: "}
	var inputs []string

	for _, value := range quoteDetails {
		util.PrintCyan(value)
		input := util.ScanUserInput()
		inputs = append(inputs, input)
	}

	return inputs
}

func ProcessUserInput(userInput []string, Database *[]Quotes) (string, bool) {

	for _, value := range userInput {
		if util.Abort(value) {
			return "<< previous action aborted by user. >>", false
		}
	}

	if FindDublicates(userInput[0], Database) {
		return "<< there are dublicates in database. >>", false
	}

	return "", true
}

func FindIndex(id int, Database *[]Quotes) int {

	index := -1

	for key, value := range *Database {
		if id == value.ID {
			index = key
		}
	}

	if index == -1 {
		panic("Error: Index not found!")
	}

	return index
}


