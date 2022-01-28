package presenter

import (
	"net/http"
	"strconv"

	"github.com/JenswBE/sunrise-alarm/src/srv-physical/api/openapi"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/entities"
)

func LedsFromEntity(color entities.PresetColor, brightness byte) openapi.Leds {
	// Convert color
	var outputColor openapi.PresetColor
	switch color {
	case entities.PresetColorBlack:
		outputColor = openapi.PRESETCOLOR_BLACK
	case entities.PresetColorRed:
		outputColor = openapi.PRESETCOLOR_RED
	case entities.PresetColorOrange:
		outputColor = openapi.PRESETCOLOR_ORANGE
	case entities.PresetColorYellow:
		outputColor = openapi.PRESETCOLOR_YELLOW
	case entities.PresetColorWarmWhite:
		outputColor = openapi.PRESETCOLOR_WARM_WHITE
	}

	// Set basic fields
	output := openapi.NewLeds(outputColor)
	output.SetBrightness(int32(brightness))
	return *output
}

func LedsToEntity(input openapi.Leds) (entities.PresetColor, byte, error) {
	// Convert color
	var color entities.PresetColor
	switch input.Color {
	case openapi.PRESETCOLOR_BLACK:
		color = entities.PresetColorBlack
	case openapi.PRESETCOLOR_RED:
		color = entities.PresetColorRed
	case openapi.PRESETCOLOR_ORANGE:
		color = entities.PresetColorOrange
	case openapi.PRESETCOLOR_YELLOW:
		color = entities.PresetColorYellow
	case openapi.PRESETCOLOR_WARM_WHITE:
		color = entities.PresetColorWarmWhite
	}

	// Validate brightness
	brightness := input.GetBrightness()
	if input.Brightness == nil {
		brightness = openapi.NewLedsWithDefaults().GetBrightness()
	}
	if brightness < 0 || brightness > 255 {
		return entities.PresetColorBlack, 0, entities.NewError(http.StatusBadRequest, openapi.ERRORCODE_UNKNOWN_ERROR, strconv.Itoa(int(brightness)), nil)
	}
	return color, byte(brightness), nil
}
