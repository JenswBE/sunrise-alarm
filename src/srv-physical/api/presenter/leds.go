package presenter

import (
	"net/http"
	"strconv"

	"github.com/JenswBE/sunrise-alarm/src/srv-physical/api/openapi"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/entities"
)

func LedsToEntity(input openapi.Leds) (entities.Color, byte, error) {
	// Convert color
	var color entities.Color
	switch input.Color {
	case openapi.PRESETCOLOR_BLACK:
		color = entities.ColorBlack
	case openapi.PRESETCOLOR_RED:
		color = entities.ColorRed
	case openapi.PRESETCOLOR_ORANGE:
		color = entities.ColorOrange
	case openapi.PRESETCOLOR_YELLOW:
		color = entities.ColorYellow
	case openapi.PRESETCOLOR_WARM_WHITE:
		color = entities.ColorWarmWhite
	}

	// Validate brightness
	brightness := input.GetBrightness()
	if input.Brightness == nil {
		brightness = openapi.NewLedsWithDefaults().GetBrightness()
	}
	if brightness < 0 || brightness > 255 {
		return entities.Color{}, 0, entities.NewError(http.StatusBadRequest, openapi.ERRORCODE_UNKNOWN_ERROR, strconv.Itoa(int(brightness)), nil)
	}
	return color, byte(brightness), nil
}
