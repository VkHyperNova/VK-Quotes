package cmd

/* Create, Read, Update, Delete  */
import (
	"fmt"
	"strings"
	"time"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func Add(Database *[]db.Quotes, DatabasePath string) bool {

	ReadCount = 0
	ErrorMsg = ""
	SuccessMsg = ""

	NewQuote := db.Quotes{
		ID:       db.FindID(Database),
		QUOTE:    util.FillEmptyInput(Quote, "Unknown"),
		AUTHOR:   util.FillEmptyInput(Author, "Unknown"),
		LANGUAGE: util.FillEmptyInput(Language, "English"),
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

func Update(id int, Database *[]db.Quotes, DatabasePath string) bool {

	ReadCount = 0
	ErrorMsg = ""
	SuccessMsg = ""

	index := db.FindIndex(id, Database)
	
	CurrentQuoteIndex = index

	(*Database)[index].QUOTE = util.FillEmptyInput(Quote, "Unknown")
	(*Database)[index].AUTHOR = util.FillEmptyInput(Author, "Unknown")
	(*Database)[index].LANGUAGE = util.FillEmptyInput(Language, "English")

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
