package db

import (
	"encoding/json"
	"strings"
	"vk-quotes/pkg/util"
)

// var DATABASE []Quotes
// var LastItemIndex = -1

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

func ReadDB(DatabasePath string) []Quotes {
	file := util.ReadFile(DatabasePath)

	Quotes := []Quotes{}

	err := json.Unmarshal(file, &Quotes)
	util.HandleError(err)

	return Quotes
}

func ValidateRequiredFiles(DatabasePath string) {
	if !util.DoesDirectoryExist(DatabasePath) {
		util.CreateDirectory("database")
		util.WriteDataToFile(DatabasePath, []byte("[]"))
		util.PrintRed("New Database Created!\n")
	}
}

// FindUniqueID finds unique id for a Quote in the database
func FindUniqueID(Database *[]Quotes) int {
	// if database is empty, return 1
	if len(*Database) == 0 {
		return 1
	}
	// return last id + 1
	return (*Database)[len(*Database)-1].ID + 1
}

func GetIndexFromId(id int, Database *[]Quotes) int {

	index := -1

	for key, value := range *Database {
		if id == value.ID {
			index = key
		}
	}

	return index
}

func CheckDublicates(quote string, Database *[]Quotes) int {
	for key, value := range *Database {
		if strings.EqualFold(value.QUOTE, quote) {
			return key
		}
	}
	return -1
}

func GetAllNames(s string, Database *[]Quotes) []string {

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

func CountNames(s string, names []string, Database *[]Quotes) map[string]int {

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
