package main

import (
	"fmt"
	"vk-quotes/pkg/cmd"
	"vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func main() {
	
	settings := util.Settings{}
	quotes := db.Quotes{}

	err := quotes.ReadFromFile()
	if err != nil {
		fmt.Println("Error loading quotes:", err)
	}
	
	cmd.CMD(&quotes, &settings)
}
