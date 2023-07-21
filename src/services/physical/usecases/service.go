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
	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories/mockdisplay"
	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories/mockleds"
	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories/mocklightsensor"
	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories/p9813leds"
	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories/rpidisplay"
	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories/tsl2591lightsensor"
	"github.com/JenswBE/sunrise-alarm/src/services/physical/utils/autobacklight"
	"github.com/JenswBE/sunrise-alarm/src/services/physical/utils/buttonpoller"
	"github.com/JenswBE/sunrise-alarm/src/services/physical/utils/buzzersequencer"
	"github.com/JenswBE/sunrise-alarm/src/utils/pubsub"
	"github.com/rs/zerolog/log"
	"github.com/stianeikeland/go-rpio/v4"
)

type PhysicalService struct {
	seq *buzzersequencer.BuzzerSequencer

	isMocked        bool
	display         repositories.Display
	leds            repositories.Leds
	lightSensor     repositories.LightSensor
	sunriseStop     chan bool
	sunriseLock     sync.Mutex // Ensures we don't mess with sunriseStop between nil check and sunriseStop creation/deletion
	sunriseDuration time.Duration
}

func NewPhysicalService(config config.PhysicalConfig, pubSub pubsub.PubSub) *PhysicalService {
	// Setup devices
	var devButton repositories.Button
	var devBuzzer repositories.Buzzer
	var devDisplay repositories.Display
	var devLeds repositories.Leds
	var devLightSensor repositories.LightSensor
	buttonChannel := make(chan buttonpoller.ButtonPress)
	if !config.Mocked {
		// Init real devices
		var err error
		if err = rpio.Open(); err != nil {
			log.Fatal().Err(err).Msg("RPIO: Failed to initialize GPIO library")
		}
		devButton = gpiobutton.NewGPIOButton(config.Button.GPIONum, true)
		devBuzzer = gpiobuzzer.NewGPIOBuzzer(config.Buzzer.GPIONum)
		devDisplay, err = rpidisplay.NewRPiDisplay()
		if err != nil {
			log.Fatal().Err(err).Msg("Display: Failed to initialize RPi display")
		}
		devLeds, err = p9813leds.NewP9813Leds()
		if err != nil {
			log.Fatal().Err(err).Msg("LED: Failed to initialize P9813 led driver on SPI0")
		}
		devLightSensor, err = tsl2591lightsensor.NewTSL2591LightSensor(config.LightSensor.I2CDevice)
		if err != nil {
			log.Fatal().Err(err).Msg("LED: Failed to initialize TSL2591 light sensor")
		}

		// Init automatic backlight adjustment
		autobacklight.NewAutoBacklight(devDisplay, devLightSensor)
	} else {
		// Init mocked devices
		devButton = mockbutton.NewMockButton()
		devBuzzer = mockbuzzer.NewMockBuzzer()
		devDisplay = mockdisplay.NewMockDisplay()
		devLeds = mockleds.NewMockLeds()
		devLightSensor = mocklightsensor.NewMockLightSensor(1234)
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
		seq:             buzzersequencer.NewBuzzerSequencer(devBuzzer, config.Buzzer.Sequence),
		isMocked:        config.Mocked,
		display:         devDisplay,
		leds:            devLeds,
		lightSensor:     devLightSensor,
		sunriseDuration: config.Leds.SunriseDuration,
	}
}

func (s *PhysicalService) Close() {
	if err := s.display.Close(); err != nil {
		log.Error().Err(err).Msg("PhysicalService.Close: Failed to close display")
	}
	s.leds.Close()
	if err := s.lightSensor.Close(); err != nil {
		log.Error().Err(err).Msg("PhysicalService.Close: Failed to close light sensor")
	}
	if !s.isMocked {
		log.Info().Msg("PhysicalService.Close: Gracefully closing physical devices")
		if err := rpio.Close(); err != nil {
			log.Error().Err(err).Msg("PhysicalService.Close: Failed to close rpio library")
		}
	}
}
