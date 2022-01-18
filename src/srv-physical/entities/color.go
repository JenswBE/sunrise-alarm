package entities

type Color struct {
	Red   byte
	Green byte
	Blue  byte
}

var (
	ColorBlack     = Color{0, 0, 0}
	ColorRed       = Color{255, 0, 0}
	ColorOrange    = Color{255, 100, 0}
	ColorYellow    = Color{255, 255, 0}
	ColorWarmWhite = Color{239, 197, 59}
)
