package cmd

import (
	"math/rand"
	"sort"
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

func PrintStatistics() {
	util.PrintCyan("\n\n\t<< Statistics >>\n")

	authors := db.SortNames("authors")
	authorsMap := db.CountNames("authors", authors)
	PrintSortedMap(authorsMap, "Authors")

	// PrintMap(authorsMap, "Authors")

	languages := db.SortNames("languages")
	languagesMap := db.CountNames("languages", languages)
	PrintSortedMap(languagesMap, "Languages")

	util.PressAnyKey()
	util.ClearScreen()
	CMD()
}

func PrintSortedMap(myMap map[string]int, name string) {

	// Make Pairs
	type pair struct {
		name  string
		count int
	}

	var pairs []pair
	for key, value := range myMap {
		pairs = append(pairs, pair{key, value})
	}

	// Sort Pairs
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	// Print Pairs
	util.PrintCyan("\n< " + name + " >\n\n")
	for _, pair := range pairs {
		util.PrintGray("[" + strconv.Itoa(pair.count) + "] ")
		util.PrintGreen(pair.name + "\n")
	}
}
