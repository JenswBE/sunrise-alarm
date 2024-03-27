package rpidisplay

import (
	"fmt"
	"os"
	"path"

	"github.com/JenswBE/sunrise-alarm/src/services/physical/repositories"
	"github.com/rs/zerolog/log"
)

var _ repositories.Display = &RPiDisplay{}

const (
	brightnessFileBasePath = "/sys/class/backlight"
	minBrightness          = 7
	maxBrightness          = 200
)

type RPiDisplay struct {
	currentBrightness  byte
	brightnessFile     *os.File
	brightnessFilePath string
}

func NewRPiDisplay() (*RPiDisplay, error) {
	// Derive brightness file path
	backlights, err := os.ReadDir(brightnessFileBasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to list backlight folders in %s: %w", brightnessFileBasePath, err)
	}
	switch len(backlights) {
	case 0:
		return nil, fmt.Errorf("no backlight folders found in %s", brightnessFileBasePath)
	case 2:
		log.Warn().Str("base_path", brightnessFileBasePath).Any("folders", backlights).Msg("Multiple backlight folders found. Using first one")
	}
	brightnessFilePath := path.Join(brightnessFileBasePath, backlights[0].Name(), "brightness")

	// Open brightness file
	log.Info().Str("path", brightnessFilePath).Msg("Opening brightness file ...")
	brightnessFile, err := os.Create(brightnessFilePath)
	if err != nil {
		log.Err(err).Str("path", brightnessFilePath).Msg("NewRPiDisplay: Failed to open brightness file")
		return nil, fmt.Errorf("failed to open brightness file: %w", err)
	}
	return &RPiDisplay{brightnessFile: brightnessFile}, nil
}

func (d *RPiDisplay) SetBrightness(brightness byte) error {
	// Setup logger
	logger := log.With().Str("path", d.brightnessFilePath).Uint8("brightness", brightness).Logger()

	// Ensure brightness is within limits
	switch {
	case brightness < minBrightness:
		brightness = minBrightness
	case brightness > maxBrightness:
		brightness = maxBrightness
	}

	// Check if value should be updated
	if brightness == d.currentBrightness {
		return nil // No update needed
	}

	// Truncate file
	if err := d.brightnessFile.Truncate(0); err != nil {
		logger.Err(err).Msg("RPiDisplay.SetBrightness: Failed to truncate brightness file")
		return fmt.Errorf("failed to truncate brightness file: %w", err)
	}

	// Write new brightness
	brightnessContent := fmt.Sprintln(brightness)
	if _, err := d.brightnessFile.WriteAt([]byte(brightnessContent), 0); err != nil {
		logger.Error().Err(err).Msg("RPiDisplay.SetBrightness: Failed to write new value to brightness file")
		return fmt.Errorf("failed to write new value to brightness file: %w", err)
	}
	d.currentBrightness = brightness
	return nil
}

func (d *RPiDisplay) Close() error {
	if err := d.brightnessFile.Close(); err != nil {
		log.Err(err).Msg("RPiDisplay.Close: Failed to close brightness file")
		return fmt.Errorf("failed to close brightness file: %w", err)
	}
	return nil
}
