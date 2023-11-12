package print

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"vk-quotes/pkg/global"
)

func PrintCommands() {
	PrintCyan("\n<< ")
	AddBrackets("add")
	AddBrackets("update")
	AddBrackets("delete")
	AddBrackets("q")
	PrintCyan(" >>\n")
}

func AddBrackets(name string) {
	PrintCyan("[")
	PrintYellow(name)
	PrintCyan("] ")
}

func PrintProgramStart() {
	PrintCyan("\n<< VK-QUOTES " + global.Version + " >>\n")
}

func ClearScreen() {

	if runtime.GOOS == "linux" {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func PrintQuote(index int) {
	PrintYellow(strconv.Itoa(global.DB[index].ID) + ". ")
	PrintYellow(global.DB[index].QUOTE + " ")
}

func PrintQuotes() {
	PrintCyan("\n\n<< Quotes >>\n")

	for _, value := range global.DB {
		PrintGray(strconv.Itoa(value.ID) + ". ")
		PrintCyan("Quote: \"")
		PrintGreen(value.QUOTE + "\"")
		PrintCyan(" By ")
		PrintCyan(value.AUTHOR + "\n")
	}
}

