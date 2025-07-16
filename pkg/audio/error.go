package audio

import (
	"embed"
	"fmt"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

//go:embed "error.mp3"
var ErrorSound embed.FS

func PlayErrorSound() error {
	// Open the embedded file
	file, err := ErrorSound.Open("error.mp3")
	if err != nil {
		return fmt.Errorf("failed to open embedded MP3: %w", err)
	}

	// Decode the MP3
	streamer, format, err := mp3.Decode(file)
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to decode MP3: %w", err)
	}

	// Initialize the speaker with the sample rate
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		streamer.Close()
		file.Close()
		return fmt.Errorf("failed to initialize speaker: %w", err)
	}

	// Play the sound
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	// Wait for playback to finish
	<-done

	// Clean up
	streamer.Close()
	file.Close()

	return nil
}
