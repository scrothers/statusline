package style

import "testing"

func TestChunkWidth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		text string
		want int
	}{
		{name: "empty", text: "", want: 0},
		{name: "ascii", text: "hello", want: 5},
		{name: "PUA glyph counts as one cell", text: "", want: 1},
		{name: "supplementary-plane glyph counts as one cell", text: "\U000F06A9", want: 1},
		{name: "mixed icon and text", text: " 2.17", want: 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := Chunk{Text: tt.text}
			if got := c.Width(); got != tt.want {
				t.Errorf("Width() = %d, want %d", got, tt.want)
			}
		})
	}
}
