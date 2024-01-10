package main

import (
	"vk-quotes/pkg/cmd"
	"vk-quotes/pkg/util"	
	db "vk-quotes/pkg/db"	
)

func main() {
	util.ClearScreen()
	db.ValidateRequiredFiles()
	cmd.CMD()
}
