package p9813led

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/stianeikeland/go-rpio/v4"
)

// var _ repositories.Button = &GPIOButton{}

type P9813LED struct{}

// NewP9813LED starts a new SPI session for P9813 LED driver on SPI0
func NewP9813LED() (*P9813LED, error) {
	err := rpio.SpiBegin(rpio.Spi0)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SPI for P9813 LED driver: %w", err)
	}
	// rpio.SpiMode(1, 0)
	rpio.SpiSpeed(10_000) // 10kHz
	return &P9813LED{}, nil
}

func (p *P9813LED) Close() {
	rpio.SpiEnd(rpio.Spi0)
}

func (p *P9813LED) SetColor(r, g, b byte) {
	// Data is 96 bits, first and last 32 bits are empty
	data := make([]byte, 12)

	// Build header
	var header byte = 0b1100_0000
	header |= p.getAntiCode(b) << 4
	header |= p.getAntiCode(g) << 2
	header |= p.getAntiCode(r)
	data[4] = header

	// Set color data
	data[5] = b
	data[6] = g
	data[7] = r

	// Send over SPI
	log.Debug().Uint8("red", r).Uint8("green", g).Uint8("blue", b).Hex("data", data).Msg("P9813: Setting new color")
	rpio.SpiTransmit(data...)
}

func (p *P9813LED) getAntiCode(data byte) byte {
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
