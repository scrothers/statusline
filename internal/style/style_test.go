package style

import (
	"strings"
	"testing"
)

func TestParseHex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      string
		want    Color
		wantErr bool
	}{
		{name: "valid lowercase", in: "#ff8019", want: RGB(0xff, 0x80, 0x19)},
		{name: "valid black", in: "#000000", want: RGB(0, 0, 0)},
		{name: "missing hash", in: "ff8019", wantErr: true},
		{name: "wrong length", in: "#fff", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseHex(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseHex(%q) error = %v, wantErr %v", tt.in, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("ParseHex(%q) = %+v, want %+v", tt.in, got, tt.want)
			}
		})
	}
}

func TestLerp(t *testing.T) {
	t.Parallel()

	green := RGB(0, 200, 0)
	red := RGB(200, 0, 0)

	tests := []struct {
		name string
		t    float64
		want Color
	}{
		{name: "t=0 returns a", t: 0, want: green},
		{name: "t=1 returns b", t: 1, want: red},
		{name: "t=0.5 is the midpoint", t: 0.5, want: RGB(100, 100, 0)},
		{name: "t<0 clamps to a", t: -5, want: green},
		{name: "t>1 clamps to b", t: 5, want: red},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := Lerp(green, red, tt.t); got != tt.want {
				t.Errorf("Lerp(green, red, %v) = %+v, want %+v", tt.t, got, tt.want)
			}
		})
	}
}

func TestPaint(t *testing.T) {
	t.Run("emits truecolor SGR codes", func(t *testing.T) {
		got := Paint("hi", RGB(255, 0, 0), Default, false)
		want := "\x1b[38;2;255;0;0mhi\x1b[0m"
		if got != want {
			t.Errorf("Paint() = %q, want %q", got, want)
		}
	})

	t.Run("no color codes when both channels default", func(t *testing.T) {
		got := Paint("hi", Default, Default, false)
		if got != "hi" {
			t.Errorf("Paint() = %q, want unmodified text", got)
		}
	})

	t.Run("NO_COLOR disables painting", func(t *testing.T) {
		t.Setenv("NO_COLOR", "1")
		got := Paint("hi", RGB(255, 0, 0), RGB(0, 0, 0), true)
		if got != "hi" {
			t.Errorf("Paint() under NO_COLOR = %q, want unmodified text", got)
		}
	})
}

func TestStrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "strips SGR color codes",
			in:   Paint("hi", RGB(1, 2, 3), RGB(4, 5, 6), true),
			want: "hi",
		},
		{
			name: "strips OSC8 hyperlink",
			in:   Hyperlink("repo", "https://example.com"),
			want: "repo",
		},
		{
			name: "plain text passes through unchanged",
			in:   "just text",
			want: "just text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := Strip(tt.in); got != tt.want {
				t.Errorf("Strip(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestHyperlink(t *testing.T) {
	t.Parallel()
	got := Hyperlink("repo", "https://example.com")
	if !strings.Contains(got, "repo") || !strings.Contains(got, "https://example.com") {
		t.Errorf("Hyperlink() = %q, missing text or url", got)
	}
	if Strip(got) != "repo" {
		t.Errorf("Strip(Hyperlink()) = %q, want repo", Strip(got))
	}
}
