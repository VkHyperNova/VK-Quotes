package main

import (
	"log"
	"vk-quotes/pkg/audio"
	"vk-quotes/pkg/cmd"
	"vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func main() {
	quotes := db.Quotes{}
	util.CreateDirectory()
	quotes.ReadFromFile()

	go func() {
		err := audio.PlayMP3()
		if err != nil {
			log.Printf("MP3 Playback Error: %v", err)
		}
	}()

	cmd.CommandLine(&quotes)
}
