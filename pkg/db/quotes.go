package db

import (
	"encoding/json"
	"os"
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
	dataBytes, err := json.MarshalIndent(DB, "", "  ")
	util.HandleError(err)
	util.HandleError(os.WriteFile(DatabasePath, dataBytes, 0644))
	return true
}

func OpenDB(DatabasePath string) []Quotes {
	file, _ := os.ReadFile(DatabasePath)

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

func FindDublicates(quote string, Database *[]Quotes) int {

	if quote == "Unknown" || quote == "" {
		return -1
	}

	for index, value := range *Database {
		if strings.EqualFold(value.QUOTE, quote) {
			return index
		}
	}

	return -1
}

func FindIndex(id int, Database *[]Quotes) int {

	index := -1

	for key, value := range *Database {
		if id == value.ID {
			index = key
		}
	}

	if index == -1 {
		panic("FindIndex(): Index not found!")
	}

	return index
}
