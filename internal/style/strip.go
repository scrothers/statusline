package style

// Strip removes ANSI CSI sequences (e.g. SGR color codes) and OSC 8
// hyperlink sequences from s, leaving only the visible text. It's used both
// for NO_COLOR-equivalent output and for golden tests that compare rendered
// text independent of color/link changes.
func Strip(s string) string {
	runes := []rune(s)
	var b []rune
	for i := 0; i < len(runes); i++ {
		if runes[i] != '\x1b' || i+1 >= len(runes) {
			b = append(b, runes[i])
			continue
		}
		switch runes[i+1] {
		case '[':
			j := i + 2
			for j < len(runes) && !isCSIFinalByte(runes[j]) {
				j++
			}
			i = j
		case ']':
			j := i + 2
			for j < len(runes) {
				if runes[j] == '\a' {
					break
				}
				if runes[j] == '\x1b' && j+1 < len(runes) && runes[j+1] == '\\' {
					j++
					break
				}
				j++
			}
			i = j
		default:
			b = append(b, runes[i])
		}
	}
	return string(b)
}

func isCSIFinalByte(r rune) bool {
	return r >= 0x40 && r <= 0x7E
}
