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

func PrintVKQUOTES(Version string, LastItemIndex int, Database *[]db.Quotes) {
	green := color.New(color.FgGreen)
	boldGreen := green.Add(color.Bold)
	boldGreen.Println("\n<< VK-QUOTES " + Version + " >>")

	util.PrintCyan("\nQuotes: " + strconv.Itoa(len(*Database)) + "\n")

	if LastItemIndex != -1 {
		PrintQuote(LastItemIndex, Database)
		LastItemIndex = -1
	} else if len(*Database) > 0 {
		PrintRandomQuote(Database)
	}

	util.PrintGray("\n")
	Commands := [6]string{"add", "update", "delete", "showall", "stats", "q"}
	for _, value := range Commands {
		util.PrintBrackets(value)
	}
	util.PrintCyan("\n=> ")
}

func PrintQuote(index int, Database *[]db.Quotes) {
	spaces := strings.Repeat(" ", 5)
	if len((*Database)[index].QUOTE)-len((*Database)[index].AUTHOR) >= 10 {
		spaces = strings.Repeat(" ", len((*Database)[index].QUOTE)-len((*Database)[index].AUTHOR))
	}

	util.PrintBlue("\n(" + strconv.Itoa((*Database)[index].ID) + ") ")
	util.PrintGray("\"" + (*Database)[index].QUOTE + "\"")
	util.PrintCyan("\n" + spaces + " By " + (*Database)[index].AUTHOR + " (" + (*Database)[index].DATE + ")\n")
}

func PrintAllQuotes(Database *[]db.Quotes) {
	util.PrintCyan("\n\n<< All Quotes >>\n")

	for key, _ := range *Database {
		PrintQuote(key, Database)
	}

}

func PrintRandomQuote(Database *[]db.Quotes) {

	randIndex := rand.Intn(len(*Database))

	for key := range *Database {
		if key == randIndex {
			PrintQuote(key, Database)
		}

	}
}

func PrintStatistics(Database *[]db.Quotes) {
	util.PrintCyan("\n\n\t<< Statistics >>\n")

	authors := db.GetAllNames("authors", Database)
	authorsMap := db.CountNames("authors", authors, Database)
	PrintSortedMap(authorsMap, "Authors")

	languages := db.GetAllNames("languages", Database)
	languagesMap := db.CountNames("languages", languages, Database)
	for name, num := range languagesMap {
		util.PrintGray("\n\n[" + strconv.Itoa(num) + "] ")
		util.PrintGreen(name + "\n")
	}

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

	for i := 0; i < 10; i++ {
		util.PrintGray("[" + strconv.Itoa(pairs[i].count) + "] ")
		util.PrintGreen(pairs[i].name + "\n")
	}
}
