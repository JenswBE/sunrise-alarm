package p9813leds

import (
	"fmt"

	"github.com/JenswBE/sunrise-alarm/src/srv-physical/entities"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/repositories"
	"github.com/rs/zerolog/log"
	"github.com/stianeikeland/go-rpio/v4"
)

var _ repositories.Leds = &P9813Leds{}

type P9813Leds struct {
	currentColor      entities.Color
	currentBrightness byte
}

// NewP9813LED starts a new SPI session for P9813 LED driver on SPI0
func NewP9813Leds() (*P9813Leds, error) {
	err := rpio.SpiBegin(rpio.Spi0)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SPI for P9813 LED driver: %w", err)
	}
	// rpio.SpiMode(1, 0)
	rpio.SpiSpeed(10_000) // 10kHz
	return &P9813Leds{}, nil
}

func (l *P9813Leds) Close() {
	rpio.SpiEnd(rpio.Spi0)
}

func (l *P9813Leds) GetColorAndBrightness() (entities.Color, byte) {
	return l.currentColor, l.currentBrightness
}

func (l *P9813Leds) SetColorAndBrightness(color entities.Color, brightness byte) {
	// Scale color to brightness
	red := scaleColor(color.Red, brightness)
	green := scaleColor(color.Green, brightness)
	blue := scaleColor(color.Blue, brightness)

	// Data is 96 bits, first and last 32 bits are empty
	data := make([]byte, 12)

	// Build header
	var header byte = 0b1100_0000
	header |= getAntiCode(red) << 4
	header |= getAntiCode(green) << 2
	header |= getAntiCode(blue)
	data[4] = header

	// Set color data
	data[5] = blue
	data[6] = green
	data[7] = red

	// Send over SPI
	log.Debug().
		Uint8("red", color.Red).
		Uint8("scaled_red", red).
		Uint8("green", color.Green).
		Uint8("scaled_green", green).
		Uint8("blue", color.Blue).
		Uint8("scaled_blue", blue).
		Uint8("brightness", brightness).
		Hex("data", data).
		Msg("P9813: Setting new color")
	rpio.SpiTransmit(data...)

	// Store new color and brightness
	l.currentColor = color
	l.currentBrightness = brightness
}

func scaleColor(color, brightness byte) byte {
	return byte(uint32(color) * uint32(brightness) / 255)
}

func getAntiCode(data byte) byte {
	var antiCode byte

	// Check if bit 7 is unset (1000 0000)
	if data&(1<<7) == 0 {
		// Set bit 1 of anti code (0000 0010)
		antiCode |= (1 << 1)
	}

	// Check if bit 6 is unset (0100 0000)
	if data&(1<<6) == 0 {
		// Set bit 0 of anti code (0000 0001)
		antiCode |= 1
	}

	return antiCode
}
