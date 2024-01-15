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

	if db.SearchIndexByID(db.LastAddID) != -1 {
		PrintQuote(db.SearchIndexByID(db.LastAddID))
		db.LastAddID = -1
	} else if len(db.DATABASE) > 0 {
		PrintRandomQuote()
	}

	util.PrintGray("\n")
	Commands := [6]string{"add", "update", "delete", "showall", "stats", "q"}
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
	spaces := strings.Repeat(" ", 5)
	if len(db.DATABASE[index].QUOTE)-len(db.DATABASE[index].AUTHOR) >= 10 {
		spaces = strings.Repeat(" ", len(db.DATABASE[index].QUOTE)-len(db.DATABASE[index].AUTHOR))
	}

	util.PrintBlue("\n(" + strconv.Itoa(db.DATABASE[index].ID) + ") ")
	util.PrintGray("\"" + db.DATABASE[index].QUOTE + "\"")
	util.PrintCyan("\n" + spaces + " By " + db.DATABASE[index].AUTHOR + " (" + db.DATABASE[index].DATE + ")\n")
}

func PrintAllQuotes() {
	util.PrintCyan("\n\n<< All Quotes >>\n")

	for key := range db.DATABASE {
		PrintQuote(key)
	}
	util.PressAnyKey()
	util.ClearScreen()
	CMD()
}

func PrintRandomQuote() {

	randIndex := rand.Intn(len(db.DATABASE))

	for key := range db.DATABASE {
		if key == randIndex {
			PrintQuote(key)
		}

	}
}

func PrintMap(m map[string]int, n string) {

	util.PrintCyan("\n< " + n + " >\n\n")
	for name, count := range m {
		util.PrintGray("[" + strconv.Itoa(count) + "] ")
		util.PrintGreen(name + "\n")
		
	}
}

func PrintStatistics() {
	util.PrintCyan("\n\n\t<< Statistics >>\n")

	authors := db.SortNames("authors")
	authorsMap := db.CountNames("authors", authors)
	PrintMap(authorsMap, "Authors")

	languages := db.SortNames("languages")
	languagesMap := db.CountNames("languages", languages)
	PrintMap(languagesMap, "Languages")

	util.PressAnyKey()
	util.ClearScreen()
	CMD()
}
