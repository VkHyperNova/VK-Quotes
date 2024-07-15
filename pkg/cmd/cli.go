package cmd

/*
All Print functions
*/

import (
	"fmt"
	"strconv"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func PrintCLI(quotes *db.Quotes, version string) {

	PrintProgramNameAndVersion(version)
	PrintProgramMessage()

	if util.ReadMode {
		SetRandomID(quotes.Size(), quotes.FindLastId())
	}

	if util.ID == -1 {
		util.ID = quotes.FindLastId()
	}

	quotes.PrintQuote(util.ID)

	PrintCommands()
}

func PrintProgramNameAndVersion(version string) {
	util.PrintGreen("\n|| ")
	util.PrintCyan("VK-QUOTES " + version)
	util.PrintGreen(" ||")
}

func PrintProgramMessage() {

	if util.Message != "" {
		length := len(util.Message) + 5
		dots := ""
		for i := 1; i < length; i++ {
			dots += "-"
		}
		util.PrintGreen("\n" + dots + "\n" + util.Message + "\n" + dots)
	}
}

func SetRandomID(dbsize int, latestId int) {
	randomIndex := util.GetRandomNumber(len(util.IDs))
	util.ID = util.IDs[randomIndex]
	PrintReadCounter(dbsize)
	util.IDs = append(util.IDs[:randomIndex], util.IDs[randomIndex+1:]...)
}

func PrintReadCounter(dbsize int) {

	util.PrintGreen("\n[" + strconv.Itoa(util.ReadCounter) + "] ")

	percentage := float64(util.ReadCounter) / float64(dbsize) * 100

	util.PrintGray(fmt.Sprintf("%.2f", percentage) + "% ")

	i := 0

	util.PrintGray("|")

	for i < util.ReadCounter {
		util.PrintGreen("-")
		i++
	}

	util.PrintGray("|")
}

func PrintCommands() {
	Commands := [7]string{"add", "update", "delete", "read", "showall", "stats", "q"}
	util.PrintCyan("\n\n")
	for _, value := range Commands {
		util.PrintGreen("|")
		util.PrintYellow(value)
		util.PrintGreen("| ")
	}
}
