package cmd

import (
	"fmt"
	"time"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func Add(quotes *db.Quotes, inputs []string, saveFilePath string) bool {
	newID := quotes.CreateId()
	quotes.Add(db.Quote{ID: newID, QUOTE: inputs[0], AUTHOR: inputs[1], LANGUAGE: inputs[2], DATE: time.Now().Format("02.01.2006")})
	quotes.SaveToFile(saveFilePath)

	util.ReadMode = false
	util.Message = fmt.Sprintf("<< %d Quote Added! >>", newID)
	return true
}

func Update(quotes *db.Quotes, inputs []string, saveFilePath string) bool {

	quotes.Update(db.Quote{ID: util.ID, QUOTE: inputs[0], AUTHOR: inputs[1], LANGUAGE: inputs[2], DATE: time.Now().Format("02.01.2006")})
	quotes.SaveToFile(saveFilePath)

	/* Set */
	util.ReadMode = false
	util.Message = fmt.Sprintf("<< %d Quote Updated! >>", util.ID)

	return true
}

func Delete(quotes *db.Quotes, saveFilePath string) bool {
	util.ReadMode = false
	quotes.Delete(util.ID)
	quotes.SaveToFile(saveFilePath)
	util.Message = fmt.Sprintf("<< %d Quote Deleted! >>", util.ID)
	return true
}
