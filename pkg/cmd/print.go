package cmd

/*
All Print functions
*/

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"

	"github.com/fatih/color"
)

func PrintCLI(Database *[]db.Quotes) {

	PrintProgramNameAndVersion()

	/* Message */
	PrintMessage()

	/* Quote */
	PrintRandomQuote(Database)
	PrintEditedQuote(Database)
	PrintLastQuoteInDatabase(Database)

	/* ReadMode */
	PrintReadCounter(Database)

	/* Commands */
	PrintCommands()
}

func PrintProgramNameAndVersion() {
	white := color.New(color.FgCyan)
	boldGreen := white.Add(color.Bold)
	util.PrintGreen("\n|| ")
	boldGreen.Print("VK-QUOTES " + IsVersion)
	util.PrintGreen(" ||")
}

func PrintCommands() {
	Commands := [7]string{"add", "update", "delete", "read", "showall", "stats", "q"}
	util.PrintCyan("\n")
	for _, value := range Commands {
		util.PrintGreen("|")
		util.PrintYellow(value)
		util.PrintGreen("| ")
	}
	util.PrintYellow("=> ")
}

func PrintMessage() {

	if IsMessage != "" {
		length := len(IsMessage) + 5
		dots := ""
		for i := 1; i < length; i++ {
			dots += "-"
		}
		util.PrintGreen("\n" + dots + "\n" + IsMessage + "\n" + dots)
	}
}

func PrintReadCounter(Database *[]db.Quotes) {

	if IsReadMode {
		IsEditedQuote = -1
		util.PrintGreen("\n[" + strconv.Itoa(IsReadCount) + "] ")

		percentage := float64(IsReadCount) / float64(len(*Database)) * 100
		util.PrintGray(fmt.Sprintf("%.2f", percentage) + "% ")

		i := 0
		util.PrintGray("|")
		for i < IsReadCount {
			util.PrintGreen("-")
			i++
		}
		util.PrintGray("|")
	}
}

func PrintQuote(index int, Database *[]db.Quotes) {
	util.PrintCyan("\n\n" + strconv.Itoa((*Database)[index].ID) + ". ")
	util.PrintCyan("\"")
	util.PrintGray((*Database)[index].QUOTE)
	util.PrintCyan("\"")
	util.PrintCyan("\n" + strings.Repeat(" ", 50) + " By " + (*Database)[index].AUTHOR + " (" + (*Database)[index].DATE + " " + (*Database)[index].LANGUAGE + ")\n")
}

func PrintAllQuotes(Database *[]db.Quotes) {
	util.PrintCyan("\n\n<< All Quotes >>\n")

	for key := range *Database {
		PrintQuote(key, Database)
	}

	IsReadCount = 0
}

func PrintRandomQuote(Database *[]db.Quotes) {

	/* Check if index exists and append to global variable 'IsUsedIndexes' */
	
	var randomIndex int
	isValid := false
	
	for !isValid {
		randomIndex = rand.Intn(len(*Database))
		if !util.ArrayContainsInt(IsUsedIndexes, randomIndex) {
			if IsReadMode {
				IsUsedIndexes = append(IsUsedIndexes, randomIndex)
			}
			isValid = true
		}
	}

	/* Find and print random quote */
	
	if len(*Database) > 0 && IsReadMode {
		for index := range *Database {
			if index == randomIndex {
				PrintQuote(index, Database)
			}
		}
	}
}

func PrintEditedQuote(Database *[]db.Quotes) {
	
	/* Print last edited quote if exists */

	if IsEditedQuote != -1 {
		PrintQuote(IsEditedQuote, Database)
	}
}

func DeleteUsedIndexes(Database *[]db.Quotes) {

	/* Empty the Slice if its full */

	if len(IsUsedIndexes) == len(*Database) {
		IsMessage = "<< You Have Read Everything! >>"
		IsUsedIndexes = IsUsedIndexes[:0]
		IsReadCount = 0
		IsReadMode = false
	}
}

func PrintLastQuoteInDatabase(Database *[]db.Quotes) {

	/* Print last quote if not in reading mode */

	if !IsReadMode && IsEditedQuote == -1 {
		PrintQuote(len(*Database)-1, Database)
	}
}

func PrintAuthors(Database *[]db.Quotes) {

	/* Get All Author Names */

	var authors []string
	for _, value := range *Database {
		if !util.ArrayContainsString(authors, value.AUTHOR) {
			authors = append(authors, value.AUTHOR)
		}
	}

	/* Count Authors By Name */

	authorsMap := make(map[string]int)
	for _, name := range authors {
		for _, value := range *Database {
			if value.AUTHOR == name {
				authorsMap[name] += 1
			}
		}
	}

	/* Make Pairs */

	type pair struct {
		name  string
		count int
	}

	var pairs []pair
	for key, value := range authorsMap {
		pairs = append(pairs, pair{key, value})
	}

	/* Sort */

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	/* Print */
	
	for i := 0; i < len(pairs) && i < 10; i++ {
		util.PrintGray("[" + strconv.Itoa(pairs[i].count) + "] ")
		util.PrintGreen(pairs[i].name + "\n")
	}
}

func PrintLanguages(Database *[]db.Quotes) {

	languages := []string{"English", "Russian", "Estonian", "German"}

	/* Count */
	languagesMap := make(map[string]int)
	for _, name := range languages {
		for _, value := range *Database {
			if value.LANGUAGE == name {
				languagesMap[name] += 1
			}
		}
	}

	/* Print */
	for name, num := range languagesMap {
		util.PrintGray("\n[" + strconv.Itoa(num) + "] " + name)
	}
}

func PrintStatistics(Database *[]db.Quotes) {
	util.PrintCyan("\n\n\t<< Statistics >>\n")
	util.PrintCyan("\n--------------------------------------------\n\n")
	PrintAuthors(Database)
	util.PrintCyan("\n--------------------------------------------\n")
	PrintLanguages(Database)
	util.PrintCyan("\n\n--------------------------------------------\n")
	IsReadCount = 0
}
