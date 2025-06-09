package audio

import (
	"testing"
)

// TestPlayMP3 simply checks that PlayMP3 returns no error.
// It does not verify actual audio playback.
func TestPlayMP3_NoError(t *testing.T) {
	err := PlayMP3()
	if err != nil {
		t.Fatalf("PlayMP3() returned error: %v", err)
	}
}
