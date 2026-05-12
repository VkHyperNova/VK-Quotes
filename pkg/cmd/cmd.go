package cmd

import (
	"fmt"
	"log"
	"os"
	"vk-quotes/pkg/audio"
	"vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func CommandLine(q *db.Quotes) {
	go func() {
		if err := audio.PlayMP3(); err != nil {
			log.Printf("MP3 Playback Error: %v", err)
		}
	}()

	for {
		q.PrintCLI()

		input := util.ReadInput()
		if input.Raw == "" {
			continue
		}

		switch input.Command {
		case "add", "a":
			if err := q.Add(); err != nil {
				audio.PlayErrorSound()
				fmt.Println(err)
				util.PressAnyKey()
			}
		case "update", "u":
			if err := q.Update(input.ID); err != nil {
				audio.PlayErrorSound()
				fmt.Println(err)
				util.PressAnyKey()
			}
		case "delete", "d":
			if err := q.Delete(input.ID); err != nil {
				audio.PlayErrorSound()
				fmt.Println(err)
				util.PressAnyKey()
			}
		case "import", "i":
			if err := q.Import(); err != nil {
				fmt.Println(err)
			}
			util.PressAnyKey()
		case "export", "e":
			if err := q.Export(); err != nil {
				fmt.Println(err)
			}
			util.PressAnyKey()
		case "unmount":
			if err := util.UnmountDrive(); err != nil {
				fmt.Println(err)
			}
			util.PressAnyKey()
		case "history", "h":
			q.History()
		case "stats":
			q.PrintStatistics()
		case "read":
			q.Read()
		case "findsimilar":
			db.FindSimilarQuotes(q)
		case "pause":
			audio.PauseMP3()
		case "resume":
			audio.ResumeMP3()
		case "random", "r":
			q.PrintRandomQuotes(input.ID)
		case "q", "quit":
			util.ClearScreen()
			os.Exit(0)
		default:
			q.Search(input.Raw)
		}
	}
}