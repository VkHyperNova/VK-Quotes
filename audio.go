package main

import (
	"embed"
	"fmt"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

//go:embed Flute.mp3 error.mp3
var sounds embed.FS

var (
	speakerOnce sync.Once
	speakerErr  error
	speakerRate beep.SampleRate

	ctrlMu sync.Mutex
	ctrl   *beep.Ctrl
)

// InitSpeaker initializes the speaker once. Every call returns the result
// of that first initialization, including a failure.
func InitSpeaker(rate beep.SampleRate) error {
	speakerOnce.Do(func() {
		speakerRate = rate
		speakerErr = speaker.Init(rate, rate.N(time.Second/10))
	})
	return speakerErr
}

func decode(name string) (beep.StreamSeekCloser, beep.Format, error) {
	f, err := sounds.Open(name)
	if err != nil {
		return nil, beep.Format{}, fmt.Errorf("failed to open embedded %s: %w", name, err)
	}
	s, format, err := mp3.Decode(f)
	if err != nil {
		f.Close()
		return nil, beep.Format{}, fmt.Errorf("failed to decode %s: %w", name, err)
	}
	return s, format, nil
}

// matchRate resamples s to the speaker's sample rate when they differ.
func matchRate(s beep.Streamer, from beep.SampleRate) beep.Streamer {
	if from != speakerRate {
		return beep.Resample(4, from, speakerRate, s)
	}
	return s
}

func PlayMP3() error {
	s, format, err := decode("Flute.mp3")
	if err != nil {
		return err
	}
	if err := InitSpeaker(format.SampleRate); err != nil {
		s.Close()
		return fmt.Errorf("failed to initialize speaker: %w", err)
	}

	c := &beep.Ctrl{Streamer: matchRate(beep.Loop(-1, s), format.SampleRate)}

	ctrlMu.Lock()
	ctrl = c
	ctrlMu.Unlock()

	speaker.Play(c)
	return nil
}

func setPaused(paused bool) bool {
	ctrlMu.Lock()
	c := ctrl
	ctrlMu.Unlock()
	if c == nil {
		return false
	}
	speaker.Lock()
	c.Paused = paused
	speaker.Unlock()
	return true
}

func PauseMP3() {
	if setPaused(true) {
		fmt.Println("Music Paused!")
	}
}

func ResumeMP3() {
	if setPaused(false) {
		fmt.Println()
	}
}

func PlayErrorSound() error {
	s, format, err := decode("error.mp3")
	if err != nil {
		return err
	}
	if err := InitSpeaker(format.SampleRate); err != nil {
		s.Close()
		return fmt.Errorf("failed to initialize speaker: %w", err)
	}

	// Closing the streamer also closes the underlying embedded file.
	speaker.Play(beep.Seq(matchRate(s, format.SampleRate), beep.Callback(func() {
		s.Close()
	})))
	return nil
}
