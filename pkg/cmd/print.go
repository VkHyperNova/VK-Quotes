package cmd

import (
	"math/rand"
	"strconv"
	"strings"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
	"github.com/fatih/color"
)

func PrintVKQUOTES(Version string) {
	green := color.New(color.FgGreen)
	boldGreen := green.Add(color.Bold)
	boldGreen.Println("\n<< VK-QUOTES " + Version + " >>")

	util.PrintCyan("\nQuotes: " + strconv.Itoa(len(db.DATABASE)) + "\n")

	if len(db.DATABASE) > 0 {
		PrintRandomQuote()
	}
	

	util.PrintGray("\n")
	Commands := [5]string{"add", "update", "delete", "showall", "q"}
	for _, value := range Commands {
		PrintBrackets(value)
	}
	util.PrintCyan("\n=> ")
}

func PrintBrackets(name string) {
	util.PrintCyan("[")
	util.PrintYellow(name)
	util.PrintCyan("] ")
}

func PrintQuote(index int) {
	util.PrintYellow(strconv.Itoa(db.DATABASE[index].ID) + ". ")
	util.PrintYellow(db.DATABASE[index].QUOTE + " ")
}

func PrintAllQuotes() {
	util.PrintCyan("\n\n<< Quotes >>\n")

	for _, value := range db.DATABASE {
		util.PrintGray(strconv.Itoa(value.ID) + ". ")
		util.PrintCyan("Quote: \"")
		util.PrintGreen(value.QUOTE + "\"")
		util.PrintCyan(" By ")
		util.PrintCyan(value.AUTHOR + "\n")
	}
	util.Prompt("\nPress any key to return to CMD\n")
}

func PrintRandomQuote() {

	randIndex := rand.Intn(len(db.DATABASE))

	for key, value := range db.DATABASE {
		if key == randIndex {
			spaces := strings.Repeat(" ", len(value.QUOTE)-len(value.AUTHOR))
			util.PrintBlue("\n(" + strconv.Itoa(value.ID) + ") ")
			util.PrintGray("\"" + value.QUOTE + "\"")
			util.PrintBlue("\n" + spaces + " By " + value.AUTHOR + " (" + value.DATE + ")\n")
		}

	}
}
