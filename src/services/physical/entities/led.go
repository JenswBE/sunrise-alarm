package entities

type LEDState struct {
	Color      PresetColor
	Brightness byte
}

func (s LEDState) IsOff() bool {
	return s.Color == PresetColorBlack || s.Brightness == 0
}

func (s LEDState) IsOn() bool {
	return !s.IsOff()
}
