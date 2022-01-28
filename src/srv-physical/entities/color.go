package entities

type PresetColor string

const (
	PresetColorBlack     PresetColor = "BLACK"
	PresetColorRed       PresetColor = "RED"
	PresetColorOrange    PresetColor = "ORANGE"
	PresetColorYellow    PresetColor = "YELLOW"
	PresetColorWarmWhite PresetColor = "WARM_WHITE"
)

func (p PresetColor) ToRGB() RGBColor {
	switch p {
	case PresetColorBlack:
		return RGBColor{0, 0, 0}
	case PresetColorRed:
		return RGBColor{255, 0, 0}
	case PresetColorOrange:
		return RGBColor{255, 100, 0}
	case PresetColorYellow:
		return RGBColor{255, 255, 0}
	case PresetColorWarmWhite:
		return RGBColor{239, 197, 59}
	}
	return RGBColor{} // Covered by exhaustive check
}

type RGBColor struct {
	Red   byte
	Green byte
	Blue  byte
}
