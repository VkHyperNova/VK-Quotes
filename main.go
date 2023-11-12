package main

import (
	"vk-quotes/pkg/cmd"
	"vk-quotes/pkg/print"
	"vk-quotes/pkg/dir"
)

func main() {
	print.ClearScreen()
	dir.ValidateRequiredFiles()
	cmd.CMD()
}
