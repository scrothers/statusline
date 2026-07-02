package style

import "fmt"

// Color is a 24-bit truecolor value, or the zero value to mean "leave the
// terminal's own default for this channel" (no escape code emitted).
type Color struct {
	R, G, B uint8
	Valid   bool
}

// Default is the terminal's own foreground/background — painting with it
// emits no color escape code for that channel.
var Default = Color{}

// RGB builds a Color from 8-bit red, green, and blue components.
func RGB(r, g, b uint8) Color {
	return Color{R: r, G: g, B: b, Valid: true}
}

// ParseHex parses a "#rrggbb" string into a Color.
func ParseHex(s string) (Color, error) {
	if len(s) != 7 || s[0] != '#' {
		return Color{}, fmt.Errorf("style: parse hex color %q: want format #rrggbb", s)
	}
	var r, g, b uint8
	if _, err := fmt.Sscanf(s, "#%02x%02x%02x", &r, &g, &b); err != nil {
		return Color{}, fmt.Errorf("style: parse hex color %q: %w", s, err)
	}
	return RGB(r, g, b), nil
}
