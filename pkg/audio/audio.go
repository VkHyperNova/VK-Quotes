package audio

import (
	"embed"
	"fmt"
	"time"
	"vk-quotes/pkg/config"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

//go:embed "Flute.mp3"
var FluteMP3 embed.FS

var ctrl *beep.Ctrl

func PlayMP3() error {
	// Open the embedded file
	file, err := FluteMP3.Open("Flute.mp3")
	if err != nil {
		return fmt.Errorf("failed to open embedded MP3: %w", err)
	}

	// Decode the MP3
	streamer, format, err := mp3.Decode(file)
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to decode MP3: %w", err)
	}


	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	ctrl = &beep.Ctrl{Streamer: streamer, Paused: false}
	speaker.Play(ctrl)

	return nil // Don't block here; let main continue
}

func PauseMP3() {
	if ctrl != nil {
		speaker.Lock()
		ctrl.Paused = true
		speaker.Unlock()
	}
	config.AddMessage(config.Yellow + "Music Paused!" + config.Reset)
}

func ResumeMP3() {
	if ctrl != nil {
		speaker.Lock()
		ctrl.Paused = false
		speaker.Unlock()
	}
	config.AddMessage(config.Green + "Music Resumed!" + config.Reset)
}

