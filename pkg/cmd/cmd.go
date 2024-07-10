package cmd

import (
	"fmt"
	"time"
	db "vk-quotes/pkg/db"
)

func LoadQuotes(filepath string) db.Quotes {
	quotes := db.Quotes{}
	err := quotes.ReadFromFile(filepath)
	if err != nil {
		fmt.Println("Error loading quotes:", err)
	}

	return quotes
}

func Add(quotes *db.Quotes, saveFilePath string) bool {
	db.ReadMode = false
	inputs, validation := UserInput(quotes)
	if validation {
		newID := quotes.CreateId()
		quotes.Add(db.Quote{ID: newID, QUOTE: inputs[0], AUTHOR: inputs[1], LANGUAGE: inputs[2], DATE: time.Now().Format("02.01.2006")})
		quotes.SaveToFile(saveFilePath)
		quotes.GetLastId()
		PrintMessage = fmt.Sprintf("<< %d Quote Added! >>", newID)
		return true
	}
	return false
}

func Update(quotes *db.Quotes, id int, saveFilePath string) bool {

	index := quotes.FindIndex(id)
	if index == -1 {
		PrintMessage = fmt.Sprintf("<< %d Index Not Found! >>", id)
		return false
	}
	
	updatedInputs := UpdateUserInput(quotes, index)
	quotes.Update(db.Quote{ID: id, QUOTE: updatedInputs[0], AUTHOR: updatedInputs[1], LANGUAGE: updatedInputs[2], DATE: time.Now().Format("02.01.2006")})
	quotes.SaveToFile(saveFilePath)

	/* Set */
	db.ReadMode = false
	PrintID = id
	PrintMessage = fmt.Sprintf("<< %d Quote Updated! >>", id)
	return true
}

func Delete(quotes *db.Quotes, id int, saveFilePath string) bool {
	db.ReadMode = false
	quotes.Delete(id)
	quotes.SaveToFile(saveFilePath)
	quotes.GetLastId()
	PrintMessage = fmt.Sprintf("<< %d Quote Deleted! >>", id)
	return true
}
