package style

// Hyperlink wraps text in an OSC 8 escape sequence linking to url. Terminals
// that don't support OSC 8 (or contexts where Strip is applied) simply show
// or keep the plain text.
func Hyperlink(text, url string) string {
	return "\x1b]8;;" + url + "\a" + text + "\x1b]8;;\a"
}
