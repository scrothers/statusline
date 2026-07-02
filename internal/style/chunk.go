package style

// Chunk is one rendered unit of a status line: a run of text with a
// foreground/background color pair, ready to be joined by the powerline
// layout logic in internal/render.
type Chunk struct {
	Text string
	FG   Color
	BG   Color
	Bold bool
}

// Width returns the chunk's display width in terminal cells, assuming every
// rune (including Nerd Font private-use-area glyphs) is single-width.
func (c Chunk) Width() int {
	width := 0
	for range c.Text {
		width++
	}
	return width
}
