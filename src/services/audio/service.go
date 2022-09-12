package audio

import (
	"github.com/JenswBE/sunrise-alarm/src/services/audio/usecases"
)

type Service interface {
	PlayMusic() error
	StopMusic() error
	IncreaseVolume() error
}

func NewAudioService() (Service, error) {
	return usecases.NewAudioService()
}
