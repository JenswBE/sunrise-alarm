package usecases

import (
	"errors"
	"fmt"

	"github.com/JenswBE/sunrise-alarm/src/services/audio/data"
	"github.com/JenswBE/sunrise-alarm/src/services/audio/utils/loop"
	"github.com/hajimehoshi/go-mp3"
)

func (s *AudioService) PlayMusic() error {
	if s.otoPlayer != nil && s.otoPlayer.IsPlaying() {
		// No action required
		return nil
	}

	// To keep things simple, currently only the default song is supported.
	// NOTE TO SELF: Once implementing support for additional songs, make sure to resample to a single sample rate.
	d, err := mp3.NewDecoder(loop.New(data.DefaultSong))
	if err != nil {
		return fmt.Errorf("failed to create new mp3 decoder: %w", err)
	}

	// Create oto player
	s.otoPlayer = s.otoContext.NewPlayer(d)
	s.otoPlayer.SetVolume(0)
	s.otoPlayer.Play()
	return nil
}

func (s *AudioService) StopMusic() error {
	if s.otoPlayer == nil {
		// No action required
		return nil
	}

	err := s.otoPlayer.Close()
	if err != nil {
		return fmt.Errorf("failed to close oto player: %w", err)
	}
	s.otoPlayer = nil
	return nil
}

func (s *AudioService) IncreaseVolume() error {
	if s.otoPlayer == nil {
		return errors.New("no music is playing")
	}
	if s.otoPlayer.Volume() == 1 {
		// Already at max volume
		return nil
	}
	s.otoPlayer.SetVolume(s.otoPlayer.Volume() + 0.01)
	return nil
}
