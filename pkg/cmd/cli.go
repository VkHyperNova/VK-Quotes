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

func PrintCLI(quotes *db.Quotes, settings *util.Settings) {

	util.ClearScreen()
	PrintProgramNameAndVersion()
	PrintProgramMessage(settings)
	PrintQuote(settings, quotes)
	PrintCommands()
}

func PrintProgramNameAndVersion() {
	util.PrintGreen("\n|| ")
	util.PrintCyan("VK-QUOTES 1.24")
	util.PrintGreen(" ||")
}

func PrintProgramMessage(settings *util.Settings) {

	if settings.Message == "" {
		settings.Message = "Hello, world!"
	}

	if settings.Message != "" {
		length := len(settings.Message) + 5
		dots := ""
		for i := 1; i < length; i++ {
			dots += "-"
		}
		util.PrintGreen("\n" + dots + "\n" + settings.Message + "\n" + dots)
	}
}

func PrintQuote(settings *util.Settings, quotes *db.Quotes) {

	if len(settings.RandomIDs) > 0 {
		util.SetRandomID(settings)
		PrintReadCounter(settings, quotes)
	} 

	if settings.ID == 0 || settings.ID == -1 {
		quotes.SetToLastID(settings)
	}
	quotes.PrintQuote(settings.ID)
}

func PrintReadCounter(settings *util.Settings, quotes *db.Quotes) {

	util.PrintGreen("\n[" + strconv.Itoa(settings.ReadCounter) + "] ")

	percentage := float64(settings.ReadCounter) / float64(quotes.Size()) * 100

	util.PrintGray(fmt.Sprintf("%.2f", percentage) + "% ")

	i := 0

	util.PrintGray("|")

	for i < settings.ReadCounter {
		util.PrintGreen("-")
		i++
	}

	util.PrintGray("|")
}

func PrintCommands() {
	Commands := [8]string{"add", "update", "delete", "read", "showall", "stats", "similar", "q"}
	util.PrintCyan("\n\n")
	for _, value := range Commands {
		util.PrintGreen("|")
		util.PrintYellow(value)
		util.PrintGreen("| ")
	}
	util.PrintCyan("\n")
}
