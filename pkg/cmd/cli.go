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

var (
	PrintID      = -1
	PrintMessage = "<< Welcome >>"
)

func PrintCLI(quotes *db.Quotes, version string) {

	PrintProgramNameAndVersion(version)
	PrintProgramMessage()

	CheckReadMode(quotes.Size(), quotes.GetLastId())
	quotes.PrintQuote(PrintID)

	PrintCommands()
}

func PrintProgramNameAndVersion(version string) {
	util.PrintGreen("\n|| ")
	util.PrintCyan("VK-QUOTES " + version)
	util.PrintGreen(" ||")
}

func PrintProgramMessage() {

	if PrintMessage != "" {
		length := len(PrintMessage) + 5
		dots := ""
		for i := 1; i < length; i++ {
			dots += "-"
		}
		util.PrintGreen("\n" + dots + "\n" + PrintMessage + "\n" + dots)
	}
}

func CheckReadMode(dbsize int, latestId int) {
	if db.ReadMode {
		randomIndex := util.GetRandomNumber(len(db.IDs))
		PrintID = db.IDs[randomIndex]
		PrintReadCounter(dbsize)
		db.IDs = append(db.IDs[:randomIndex], db.IDs[randomIndex+1:]...)
	} else {
		PrintID = latestId
	}
}

func PrintReadCounter(dbsize int) {

	util.PrintGreen("\n[" + strconv.Itoa(db.ReadCounter) + "] ")

	percentage := float64(db.ReadCounter) / float64(dbsize) * 100

	util.PrintGray(fmt.Sprintf("%.2f", percentage) + "% ")

	i := 0

	util.PrintGray("|")

	for i < db.ReadCounter {
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
	util.PrintYellow("=> ")
}
