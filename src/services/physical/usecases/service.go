package usecases

import (
	"sync"
	"time"

	"github.com/JenswBE/sunrise-alarm/src/cmd/config"
	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories"
	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories/gpiobutton"
	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories/gpiobuzzer"
	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories/mockbutton"
	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories/mockbuzzer"
	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories/mockleds"
	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories/p9813leds"
	"github.com/JenswBE/sunrise-alarm/src/services/physical/utils/buttonpoller"
	"github.com/JenswBE/sunrise-alarm/src/services/physical/utils/buzzersequencer"
	"github.com/JenswBE/sunrise-alarm/src/utils/pubsub"
	"github.com/rs/zerolog/log"
	"github.com/stianeikeland/go-rpio/v4"
)

type PhysicalService struct {
	seq *buzzersequencer.BuzzerSequencer

	isMocked        bool
	leds            repositories.Leds
	sunriseStop     chan bool
	sunriseLock     sync.Mutex // Ensures we don't mess with sunriseStop between nil check and sunriseStop creation/deletion
	sunriseDuration time.Duration
}

func NewPhysicalService(config config.PhysicalConfig, pubSub pubsub.PubSub) *PhysicalService {
	// Setup devices
	var devButton repositories.Button
	var devBuzzer repositories.Buzzer
	var devLeds repositories.Leds
	buttonChannel := make(chan buttonpoller.ButtonPress)
	if !config.Mocked {
		// Init real devices
		if err := rpio.Open(); err != nil {
			log.Fatal().Err(err).Msg("RPIO: Failed to initialize GPIO library")
		}
		devButton = gpiobutton.NewGPIOButton(config.Button.GPIONum, true)
		devBuzzer = gpiobuzzer.NewGPIOBuzzer(config.Buzzer.GPIONum)
		p9813Leds, err := p9813leds.NewP9813Leds()
		if err != nil {
			log.Fatal().Err(err).Msg("LED: Failed to initialize P9813 led driver on SPI0")
		}
		devLeds = p9813Leds
	} else {
		// Init mocked devices
		devButton = mockbutton.NewMockButton()
		devBuzzer = mockbuzzer.NewMockBuzzer()
		devLeds = mockleds.NewMockLeds()
	}

	// Start polling button
	buttonpoller.NewButtonPoller(devButton, buttonChannel)
	go func() {
		for {
			switch <-buttonChannel {
			case buttonpoller.ButtonPressShort:
				pubSub.Publish(pubsub.EventButtonPressedShort{})
			case buttonpoller.ButtonPressLong:
				pubSub.Publish(pubsub.EventButtonPressedLong{})
			}
		}
	}()

	// Build service
	return &PhysicalService{
		seq:             buzzersequencer.NewBuzzerSequencer(devBuzzer),
		isMocked:        config.Mocked,
		leds:            devLeds,
		sunriseDuration: config.Leds.SunriseDuration,
	}
}

func (s *PhysicalService) Close() {
	if !s.isMocked {
		log.Info().Msg("PhysicalService.Close: Gracefully closing physical devices")
		if err := rpio.Close(); err != nil {
			log.Error().Err(err).Msg("PhysicalService.Close: Failed to close rpio library")
		}
		if leds, ok := s.leds.(*p9813leds.P9813Leds); ok {
			leds.Close()
		}
	}
}
