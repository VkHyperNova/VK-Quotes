package db

import (
	"encoding/json"
	"strings"
	"vk-quotes/pkg/util"
)

var DATABASE []Quotes
var LastItemID = -1

type Quotes struct {
	ID       int    `json:"id"`
	QUOTE    string `json:"quote"`
	AUTHOR   string `json:"author"`
	LANGUAGE string `json:"language"`
	DATE     string `json:"date"`
}

func SaveDB(DatabasePath string) string {
	DatabaseAsByte := util.InterfaceToByte(DATABASE)
	util.WriteDataToFile(DatabasePath, DatabaseAsByte)
	return "Database Updated!"
}

func LoadDB(DatabasePath string) []Quotes {
	file := util.ReadFile(DatabasePath)
	data := GetQuotesArray(file)

	return data
}

func ValidateRequiredFiles(DatabasePath string) {
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

func GetQuotesArray(body []byte) []Quotes {

	QuotesStruct := []Quotes{}

	err := json.Unmarshal(body, &QuotesStruct)
	util.HandleError(err)

	return QuotesStruct
}

func SearchIndexByID(id int) int {

	index := -1

	for key, value := range DATABASE {
		if id == value.ID {
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

func GetAllNames(s string) []string {

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
