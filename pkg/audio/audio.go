package audio

import (
	"embed"
	"fmt"
	"sync"
	"time"
	"vk-quotes/pkg/config"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

//go:embed "Flute.mp3"
var FluteMP3 embed.FS

var ctrl *beep.Ctrl

var speakerOnce sync.Once

func InitSpeaker(sampleRate beep.SampleRate) error {
	var err error
	speakerOnce.Do(func() {
		err = speaker.Init(sampleRate, sampleRate.N(time.Second/10))
	})
	return err
}

func PlayMP3() error {
	file, err := FluteMP3.Open("Flute.mp3")
	if err != nil {
		return fmt.Errorf("failed to open embedded MP3: %w", err)
	}

	streamer, format, err := mp3.Decode(file)
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to decode MP3: %w", err)
	}

	err = InitSpeaker(format.SampleRate)
	if err != nil {
		return fmt.Errorf("failed to initialize speaker: %w", err)
	}

	ctrl = &beep.Ctrl{Streamer: streamer, Paused: false}
	speaker.Play(ctrl)

	return nil
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

//go:embed "error.mp3"
var ErrorSound embed.FS

func PlayErrorSound() error {
	file, err := ErrorSound.Open("error.mp3")
	if err != nil {
		return fmt.Errorf("failed to open embedded MP3: %w", err)
	}

	streamer, format, err := mp3.Decode(file)
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to decode MP3: %w", err)
	}

	err = InitSpeaker(format.SampleRate)
	if err != nil {
		return fmt.Errorf("failed to initialize speaker: %w", err)
	}

	s := streamer
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		s.Close()
		file.Close()
	})))

	return nil
}


