package cmd

/* Create, Read, Update, Delete  */
import (
	"fmt"
	"strings"
	"time"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func Create(inputs []string, Database *[]db.Quotes, DatabasePath string) bool {

	ErrorMsg = ""
	SuccessMsg = ""

	NewQuote := db.Quotes{
		ID:       db.FindID(Database),
		QUOTE:    util.FillEmptyInput(inputs[0], "Unknown"),
		AUTHOR:   util.FillEmptyInput(inputs[1], "Unknown"),
		LANGUAGE: util.FillEmptyInput(inputs[2], "English"),
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

func Update(id int, inputs []string, Database *[]db.Quotes, DatabasePath string) bool {

	ErrorMsg = ""
	SuccessMsg = ""

	index := db.FindIndex(id, Database)

	CurrentQuoteIndex = index

	(*Database)[index].QUOTE = util.FillEmptyInput(inputs[0], (*Database)[index].QUOTE)
	(*Database)[index].AUTHOR = util.FillEmptyInput(inputs[1], (*Database)[index].AUTHOR)
	(*Database)[index].LANGUAGE = util.FillEmptyInput(inputs[2], (*Database)[index].LANGUAGE)

	db.SaveDB(Database, DatabasePath)

	SuccessMsg = fmt.Sprintf("<< %d Quote Updated! >>", (*Database)[index].ID)

	return true
}

func Delete(id int, Database *[]db.Quotes, DatabasePath string) bool {

	ErrorMsg = ""
	SuccessMsg = ""

	index := db.FindIndex(id, Database)

	SuccessMsg = fmt.Sprintf("<< %d Quote Deleted! >>", (*Database)[index].ID)

	(*Database) = append((*Database)[:index], (*Database)[index+1:]...)

	db.SaveDB(Database, DatabasePath)

	

	return true
}
