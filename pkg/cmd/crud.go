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

	IsReadMode = false

	NewQuote := db.Quotes{
		ID:       db.FindID(Database),
		QUOTE:    util.FillEmptyInput(IsQuote, "Unknown"),
		AUTHOR:   util.FillEmptyInput(IsAuthor, "Unknown"),
		LANGUAGE: util.FillEmptyInput(IsLanguage, "English"),
		DATE:     time.Now().Format("02.01.2006"),
	}

	*Database = append(*Database, NewQuote)

	db.SaveDB(Database, DatabasePath)

	index := db.FindIndex(NewQuote.ID, Database)

	IsMessage = fmt.Sprintf("<< %d Quote Added! >>", (*Database)[index].ID)

	IsEditedQuote = index

	return true
}

func FindByAuthor(Database *[]db.Quotes, searchString string) {

	IsReadMode = false

	for key, value := range *Database {
		if strings.Contains(strings.ToUpper(value.AUTHOR), strings.ToUpper(searchString)) {
			PrintQuote(key, Database)
		}
	}
}

func Update(id int, Database *[]db.Quotes, DatabasePath string) bool {
	
	IsReadMode = false

	index := db.FindIndex(id, Database)

	IsEditedQuote = index

	(*Database)[index].QUOTE = util.FillEmptyInput(IsQuote, "Unknown")
	(*Database)[index].AUTHOR = util.FillEmptyInput(IsAuthor, "Unknown")
	(*Database)[index].LANGUAGE = util.FillEmptyInput(IsLanguage, "English")

	db.SaveDB(Database, DatabasePath)

	IsMessage = fmt.Sprintf("<< %d Quote Updated! >>", (*Database)[index].ID)

	return true
}

func Delete(id int, Database *[]db.Quotes, DatabasePath string) bool {
	
	IsReadMode = false

	index := db.FindIndex(id, Database)

	IsMessage = fmt.Sprintf("<< %d Quote Deleted! >>", (*Database)[index].ID)

	(*Database) = append((*Database)[:index], (*Database)[index+1:]...)

	db.SaveDB(Database, DatabasePath)

	return true
}

func Read() {
	IsReadMode = true
	IsMessage = "<< Reading >>"
	IsEditedQuote = -1
}