package main

import (
	"fmt"
	"vk-quotes/pkg/cmd"
	"vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func main() {
	/* Problem with paths. */
	settings := util.Settings{}
	settings.Version = "1.24"
	settings.SaveQuotesPath = "./QuotesDB/quotes.json"
	settings.SaveSimilarPath = "./QuotesDB/similar.json"
	settings.SaveFolderPath = "QuotesDB"
	settings.Message = "Hello, World!"

	quotes := db.Quotes{}
	err := quotes.ReadFromFile(&settings)
	if err != nil {	
		fmt.Println("Error loading quotes:", err)
	}
	
	cmd.CMD(&quotes, &settings)
}
