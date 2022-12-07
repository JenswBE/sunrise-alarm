package rpidisplay

import (
	"fmt"
	"os"

	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories"
	"github.com/rs/zerolog/log"
)

var _ repositories.Display = &RPiDisplay{}

const brightnessFilePath = "/sys/class/backlight/10-0045/brightness"
const minBrightness = 7
const maxBrightness = 100

type RPiDisplay struct {
	brightnessFile *os.File
}

func NewRPiDisplay() (*RPiDisplay, error) {
	brightnessFile, err := os.Create(brightnessFilePath)
	if err != nil {
		log.Err(err).Str("path", brightnessFilePath).Msg("NewRPiDisplay: Failed to open brightness file")
		return nil, fmt.Errorf("failed to open brightness file: %w", err)
	}
	return &RPiDisplay{brightnessFile: brightnessFile}, nil
}

func (d *RPiDisplay) SetBrightness(brightness byte) error {
	// Setup logger
	logger := log.With().Str("path", brightnessFilePath).Uint8("brightness", brightness).Logger()

	// Truncate file
	if err := d.brightnessFile.Truncate(0); err != nil {
		logger.Err(err).Msg("RPiDisplay.SetBrightness: Failed to truncate brightness file")
		return fmt.Errorf("failed to truncate brightness file: %w", err)
	}

	// Ensure brightness is within limits
	switch {
	case brightness < minBrightness:
		brightness = minBrightness
	case brightness > maxBrightness:
		brightness = maxBrightness
	}

	// Write new brightness
	brightnessContent := fmt.Sprintln(brightness)
	if _, err := d.brightnessFile.WriteAt([]byte(brightnessContent), 0); err != nil {
		logger.Error().Err(err).Msg("RPiDisplay.SetBrightness: Failed to write new value to brightness file")
		return fmt.Errorf("failed to write new value to brightness file: %w", err)
	}
	return nil
}

func (d *RPiDisplay) Close() error {
	if err := d.brightnessFile.Close(); err != nil {
		log.Err(err).Msg("RPiDisplay.Close: Failed to close brightness file")
		return fmt.Errorf("failed to close brightness file: %w", err)
	}
	return nil
}
