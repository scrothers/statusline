package render

import (
	"os"
	"strconv"
)

// defaultColumns is used when COLUMNS is unset or unparseable — older
// Claude Code versions (before v2.1.153) and manual/test invocations don't
// set it.
const defaultColumns = 80

// Columns reads the terminal width Claude Code sets in the COLUMNS
// environment variable before invoking the statusline command. It reads the
// env var directly (rather than tput or a terminal ioctl) because Claude
// Code captures the command's stdout instead of connecting it to a real
// terminal, so those wouldn't see the actual width.
func Columns() int {
	if v := os.Getenv("COLUMNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return defaultColumns
}
