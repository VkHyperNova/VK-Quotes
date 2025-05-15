package audio

import (
	"embed"
	"fmt"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

//go:embed babel.mp3
var babelMP3 embed.FS

func PlayMP3() error {
	// Open the embedded file
	file, err := babelMP3.Open("babel.mp3")
	if err != nil {
		return fmt.Errorf("failed to open embedded MP3: %w", err)
	}

	// Decode the MP3
	streamer, format, err := mp3.Decode(file)
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to decode MP3: %w", err)
	}

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		streamer.Close()
		file.Close()
		return fmt.Errorf("failed to initialize speaker: %w", err)
	}

	// Loop the streamer infinitely
	looped := beep.Loop(-1, streamer)

	// Play the music
	speaker.Play(looped)

	// Block forever
	select {}
}


