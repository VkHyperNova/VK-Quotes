package print

import (
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"vk-quotes/pkg/global"

	"github.com/fatih/color" // bold tekst
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
	green := color.New(color.FgGreen)
	boldGreen := green.Add(color.Bold)
	boldGreen.Println("\n<< VK-QUOTES " + global.Version + " >>")
	// PrintCyan("\n<< VK-QUOTES " + global.Version + " >>\n")
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

func PrintAllQuotes() {
	PrintCyan("\n\n<< Quotes >>\n")

	for _, value := range global.DB {
		PrintGray(strconv.Itoa(value.ID) + ". ")
		PrintCyan("Quote: \"")
		PrintGreen(value.QUOTE + "\"")
		PrintCyan(" By ")
		PrintCyan(value.AUTHOR + "\n")
	}
}

func PrintRandomQuote() {
	randIndex := rand.Intn(len(global.DB))

	for key, value := range global.DB {
		if key == randIndex {
			spaces := strings.Repeat(" ", len(value.QUOTE) - len(value.AUTHOR))
			PrintGray("\n\"" + value.QUOTE + "\"")
			PrintBlue("\n" + spaces + " By " + value.AUTHOR + "\n")
		}

	}
}
