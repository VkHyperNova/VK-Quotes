package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"vk-quotes/pkg/global"
	"vk-quotes/pkg/print"
)

func GetInput(inputName string) string {
	print.PrintCyan(inputName)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	NewTaskString := scanner.Text()
	NewTaskString = strings.TrimSpace(NewTaskString)

	return NewTaskString
}

func InterfaceToByte(data interface{}) []byte {
	dataBytes, err := json.MarshalIndent(data, "", "  ")
	print.HandleError(err)

	return dataBytes
}

func GetFormattedDate() string {
	return time.Now().Format("02.01.2006")
}

func FindUniqueID() int {

	if len(global.DB) == 0 {
		return 1
	}

	return global.DB[len(global.DB)-1].ID + 1
}

func CompileQuote(quote string, author string, language string) global.Quotes {

	return global.Quotes{
		ID:       FindUniqueID(),
		QUOTE:    quote,
		AUTHOR:   author,
		LANGUAGE: language,
		DATE:     GetFormattedDate(),
	}
}

func GetQuotesArray(body []byte) []global.Quotes {

	QuotesStruct := []global.Quotes{}

	err := json.Unmarshal(body, &QuotesStruct)
	print.HandleError(err)

	return QuotesStruct
}

func SearchIndexByID(id int) int {

	index := -1

	for key, website := range global.DB {
		if id == website.ID {
			index = key
		}
	}

	return index
}

func Confirm() bool {

	user_input := Prompt("\n\nThis One?: ")

	if user_input == "n" || user_input == "no" {
		return false
	}
	return true
}

func Prompt(Question string) string {

	print.PrintCyan(Question)

	var user_input_string string
	fmt.Scanln(&user_input_string)

	return user_input_string
}
