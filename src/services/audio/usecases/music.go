package usecases

import (
	"fmt"

	"github.com/JenswBE/sunrise-alarm/src/services/audio/data"
	"github.com/JenswBE/sunrise-alarm/src/services/audio/utils/loop"
	"github.com/hajimehoshi/go-mp3"
	"github.com/rs/zerolog/log"
)

func (s *AudioService) PlayMusic() error {
	log.Debug().Msg("Audio.PlayMusic: Requested to play music")
	if s.otoPlayer != nil && s.otoPlayer.IsPlaying() {
		// No action required
		log.Info().Msg("Audio.PlayMusic: Requested to play music, but already playing. Ignoring request.")
		return nil
	}

	// To keep things simple, currently only the default song is supported.
	// NOTE TO SELF: Once implementing support for additional songs, make sure to resample to a single sample rate.
	d, err := mp3.NewDecoder(loop.New(data.DefaultSong))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create new mp3 decoder for default song")
		return fmt.Errorf("failed to create new mp3 decoder: %w", err)
	}

	// Create oto player
	s.otoPlayer = s.otoContext.NewPlayer(d)
	s.otoPlayer.SetVolume(0)
	s.otoPlayer.Play()
	log.Debug().Msg("Audio.PlayMusic: Started playing music")
	return nil
}

func (s *AudioService) StopMusic() error {
	log.Debug().Msg("Audio.StopMusic: Requested to stop music")
	if s.otoPlayer == nil {
		// No action required
		log.Info().Msg("Audio.StopMusic: Requested to stop music, but not playing. Ignoring request.")
		return nil
	}

	err := s.otoPlayer.Close()
	if err != nil {
		return fmt.Errorf("failed to close oto player: %w", err)
	}
	s.otoPlayer = nil
	log.Debug().Msg("Audio.PlayMusic: Stopped playing music")
	return nil
}

func (s *AudioService) IncreaseVolume() error {
	log.Debug().Msg("Audio.IncreaseVolume: Requested to increase volume")
	if s.otoPlayer == nil {
		log.Info().Msg("Audio.IncreaseVolume: Requested to increase volume, but not playing. Ignoring request.")
		return nil
	}
	if s.otoPlayer.Volume() == 1 {
		// Already at max volume
		log.Info().Msg("Audio.IncreaseVolume: Requested to increase volume, but already at max volume. Ignoring request.")
		return nil
	}
	s.otoPlayer.SetVolume(s.otoPlayer.Volume() + 0.01)
	log.Debug().Msg("Audio.IncreaseVolume: Volume increased")
	return nil
}
