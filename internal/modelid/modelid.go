package modelid

import (
	"regexp"
	"strings"
)

// familyNames maps a lowercase family word to its canonical display casing.
var familyNames = map[string]string{
	"opus":   "Opus",
	"sonnet": "Sonnet",
	"haiku":  "Haiku",
	"fable":  "Fable",
	"mythos": "Mythos",
}

var (
	// reVTag strips a Bedrock inference-profile version tag, e.g. the "-v2"
	// in "claude-3-5-sonnet-20241022-v2".
	reVTag = regexp.MustCompile(`-v\d+$`)

	// re1M strips a "[1m]" 1M-context marker, however it's cased.
	re1M = regexp.MustCompile(`(?i)\[1m\]`)

	// reDateToken matches an 8-digit snapshot date stamp, e.g. "20241022".
	reDateToken = regexp.MustCompile(`^\d{8}$`)

	// reVersionToken matches a 1-2 digit version component, e.g. "4" or "8".
	reVersionToken = regexp.MustCompile(`^\d{1,2}$`)
)

// Label returns the best available display label for a model: a decoded,
// normalized "<Family> <Version>" form of id when id is recognizably a
// Claude model, otherwise displayName verbatim, otherwise "".
func Label(id, displayName string) string {
	if label, ok := Decode(id); ok {
		return label
	}
	return displayName
}

// Decode normalizes a raw Claude model identifier from any known source —
// first-party, Bedrock, Vertex, OpenRouter, or a corporate gateway that
// reorders the family/version segments — into a "<Family> <Version>" label.
// It reports ok=false when id doesn't look like a Claude model at all, so
// callers can fall back to another source.
func Decode(id string) (label string, ok bool) {
	s := strings.TrimSpace(id)
	if s == "" {
		return "", false
	}

	// Strip an OpenRouter-style "provider/" prefix.
	if i := strings.LastIndex(s, "/"); i >= 0 {
		s = s[i+1:]
	}

	// Strip a Bedrock region + provider prefix: "us.anthropic." or
	// "anthropic." both collapse to nothing.
	if i := strings.Index(s, "anthropic."); i >= 0 {
		s = s[i+len("anthropic."):]
	}

	// Strip an OpenRouter variant tag, e.g. "claude-3.5-sonnet:beta".
	if i := strings.Index(s, ":"); i >= 0 {
		s = s[:i]
	}

	// Strip a Bedrock inference-profile version tag, e.g. "-v2".
	s = reVTag.ReplaceAllString(s, "")

	// Strip a Vertex dated-snapshot suffix, e.g. "@20251101".
	if i := strings.Index(s, "@"); i >= 0 {
		s = s[:i]
	}

	// The 1M-context marker carries no family/version information.
	s = re1M.ReplaceAllString(s, "")

	// OpenRouter spells versions with dots ("3.5"); normalize to the
	// dash-delimited scheme every other source uses.
	s = strings.ReplaceAll(s, ".", "-")

	var family, otherFamily string
	var versionParts []string
	sawClaude := false

	for tok := range strings.SplitSeq(s, "-") {
		if tok == "" {
			continue
		}
		lower := strings.ToLower(tok)
		switch {
		case lower == "claude":
			sawClaude = true
		case familyNames[lower] != "":
			family = familyNames[lower]
		case reDateToken.MatchString(tok):
			// Snapshot date stamp — not part of the displayed version.
		case reVersionToken.MatchString(tok):
			versionParts = append(versionParts, tok)
		case otherFamily == "":
			otherFamily = strings.ToUpper(lower[:1]) + lower[1:]
		}
	}

	if family == "" {
		if !sawClaude {
			return "", false
		}
		if otherFamily != "" {
			family = otherFamily
		} else {
			family = "Claude"
		}
	}

	label = family
	if len(versionParts) > 0 {
		label += " " + strings.Join(versionParts, ".")
	}
	return label, true
}
