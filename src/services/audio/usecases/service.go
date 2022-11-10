package usecases

import (
	"fmt"

	"github.com/hajimehoshi/oto/v2"
)

type AudioService struct {
	otoContext *oto.Context
	otoPlayer  oto.Player
}

func NewAudioService() (*AudioService, error) {
	// Create oto context
	c, ready, err := oto.NewContext(44100, 2, 2)
	if err != nil {
		return nil, fmt.Errorf("failed to create oto context: %w", err)
	}
	<-ready

	// Build service
	return &AudioService{otoContext: c}, nil
}
