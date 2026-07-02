package style

import (
	"os"
	"strconv"
	"strings"
)

// Paint wraps text in truecolor ANSI escape codes for fg/bg (each only if
// Valid) and bold, resetting afterward. When the NO_COLOR environment
// variable is set to any non-empty value, Paint returns text unchanged —
// the same code path Strip's tests rely on for a plain-text baseline.
func Paint(text string, fg, bg Color, bold bool) string {
	if os.Getenv("NO_COLOR") != "" {
		return text
	}

	var codes []string
	if bold {
		codes = append(codes, "1")
	}
	if fg.Valid {
		codes = append(codes, sgrTrueColor(38, fg))
	}
	if bg.Valid {
		codes = append(codes, sgrTrueColor(48, bg))
	}
	if len(codes) == 0 {
		return text
	}
	return "\x1b[" + strings.Join(codes, ";") + "m" + text + "\x1b[0m"
}

func sgrTrueColor(base int, c Color) string {
	return strconv.Itoa(base) + ";2;" + strconv.Itoa(int(c.R)) + ";" + strconv.Itoa(int(c.G)) + ";" + strconv.Itoa(int(c.B))
}
