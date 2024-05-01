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

func PrintCLI(Version string, CurrentQuoteIndex int, Database *[]db.Quotes) {

	/* Print Program Name and Version in bold */
	white := color.New(color.FgCyan)
	boldGreen := white.Add(color.Bold)
	util.PrintGreen("\n|| ")
	boldGreen.Print("VK-QUOTES " + Version)
	util.PrintGreen(" ||")

	/* Print Read Counter */
	util.PrintGreen("\n\n[" + strconv.Itoa(ReadCount) + "] ")

	/* Print Loading bar */
	percentage := float64(ReadCount) / float64(len(*Database)) * 100
	util.PrintGray(fmt.Sprintf("%.2f", percentage) + "% ")

	i := 0
	util.PrintGray("|")
	for i < ReadCount {
		util.PrintGreen("-")
		i++
	}
	util.PrintGray("|")

	/* Print Main Quote */
	if CurrentQuoteIndex != -1 {
		PrintQuote(CurrentQuoteIndex, Database)
	} else if len(*Database) > 0 {
		PrintRandomQuote(Database)
	}

	/* Print Message */
	if SuccessMsg != "" {
		util.PrintGreen("\n" + SuccessMsg + "\n")
		SuccessMsg = ""
	}

	if ErrorMsg != "" {
		util.PrintRed("\n" + ErrorMsg + "\n")
		ErrorMsg = ""
	}

	/* Print CLI */
	Commands := [6]string{"add", "update", "delete", "showall", "stats", "q"}
	util.PrintCyan("\n")
	for _, value := range Commands {
		util.PrintGreen("|")
		util.PrintYellow(value)
		util.PrintGreen("| ")
	}
	util.PrintGreen("\n|")
	util.PrintYellow("=> ")
}

func PrintQuote(index int, Database *[]db.Quotes) {
	spaces := strings.Repeat(" ", 5)

	quoteLength := len((*Database)[index].QUOTE)
	authorLength := len((*Database)[index].AUTHOR)
	repeatTimes := quoteLength - authorLength

	if repeatTimes >= 10 {
		spaces = strings.Repeat(" ", repeatTimes)
	}

	util.PrintCyan("\n\n" + strconv.Itoa((*Database)[index].ID) + " ")
	util.PrintGray("\"" + (*Database)[index].QUOTE + "\"")
	util.PrintCyan("\n" + spaces + " By " + (*Database)[index].AUTHOR + " (" + (*Database)[index].DATE + ")\n")
}

func PrintAllQuotes(Database *[]db.Quotes) {
	util.PrintCyan("\n\n<< All Quotes >>\n")

	for key := range *Database {
		PrintQuote(key, Database)
	}
}

func PrintRandomQuote(Database *[]db.Quotes) {

	var randomIndex int
	/* Validate Index */
	isValid := false
	for !isValid {
		randomIndex = rand.Intn(len(*Database))
		if !util.ArrayContainsInt(UsedIndexes, randomIndex) {
			UsedIndexes = append(UsedIndexes, randomIndex)
			isValid = true
		}
	}

	/* Find and Print Quote */
	for index := range *Database {
		if index == randomIndex {
			PrintQuote(index, Database)
		}

	}

	/* Empty the Slice if its full */
	if len(UsedIndexes) == len(*Database) {
		SuccessMsg = "You Have Read Everything!"
		UsedIndexes = UsedIndexes[:0]
		ReadCount = 0
	}
}

func PrintAuthors(Database *[]db.Quotes) {

	/* Get All Author Names*/
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
		util.PrintGray("\n[" + strconv.Itoa(num) + "] ")
		util.PrintGray(name)
	}
}

func PrintStatistics(Database *[]db.Quotes) {
	util.PrintCyan("\n\n\t<< Statistics >>\n")
	util.PrintCyan("\n--------------------------------------------\n\n")
	PrintAuthors(Database)
	util.PrintCyan("\n--------------------------------------------\n")
	PrintLanguages(Database)
	util.PrintCyan("\n\n--------------------------------------------\n")
}
