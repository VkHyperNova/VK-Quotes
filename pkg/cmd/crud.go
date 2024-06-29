package cmd

/* Create, Read, Update, Delete  */
import (
	"fmt"
	"strings"
	"time"
	db "vk-quotes/pkg/db" 
)

func Create(inputs []string, Database *[]db.Quotes, DatabasePath string) bool {

	ReadCount = 0
	ErrorMsg = ""
	SuccessMsg = ""

	NewQuote := db.Quotes{
		ID:       db.FindID(Database),
		QUOTE:    inputs[0],
		AUTHOR:   inputs[1],
		LANGUAGE: inputs[2],
		DATE:     time.Now().Format("02.01.2006"),
	}

	*Database = append(*Database, NewQuote)

	db.SaveDB(Database, DatabasePath)

	index := db.FindIndex(NewQuote.ID, Database)

	SuccessMsg = fmt.Sprintf("<< %d Quote Added! >>", (*Database)[index].ID)

	CurrentQuoteIndex = index

	return true
}

func FindByAuthor(Database *[]db.Quotes, searchString string) {

	for key, value := range *Database {
		if strings.Contains(strings.ToUpper(value.AUTHOR), strings.ToUpper(searchString)) {
			PrintQuote(key, Database)
		}
	}
}

func Update(id int, input []string, Database *[]db.Quotes, DatabasePath string) bool {

	ReadCount = 0
	ErrorMsg = ""
	SuccessMsg = ""

	index := db.FindIndex(id, Database)

	CurrentQuoteIndex = index

	(*Database)[index].QUOTE = input[0]
	(*Database)[index].AUTHOR = input[1]
	(*Database)[index].LANGUAGE = input[2]

	db.SaveDB(Database, DatabasePath)

	SuccessMsg = fmt.Sprintf("<< %d Quote Updated! >>", (*Database)[index].ID)

	return true
}

func Delete(id int, Database *[]db.Quotes, DatabasePath string) bool {

	CurrentQuoteIndex = -1
	ReadCount = 0
	ErrorMsg = ""
	SuccessMsg = ""

	index := db.FindIndex(id, Database)

	SuccessMsg = fmt.Sprintf("<< %d Quote Deleted! >>", (*Database)[index].ID)

	(*Database) = append((*Database)[:index], (*Database)[index+1:]...)

	db.SaveDB(Database, DatabasePath)

	

	return true
}
