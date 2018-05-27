package core

type Color uint8

const (
	Black Color = 0
	White Color = 1
)

func MakeColor(i uint8) Color {
	if i == 0 {
		return Black
	}
	if i == 1 {
		return White
	}
	panic("Cannot make color for given value")
}

func (c Color) Name() string {
	if c == White {
		return "white"
	}
	if c == Black {
		return "black"
	}

	panic("Cannot get name for color")
}
