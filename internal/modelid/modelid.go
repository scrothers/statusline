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
	// reSeparators normalizes underscores and whitespace — used by some
	// gateways and env-derived IDs instead of dashes — to the dash-delimited
	// scheme every other source uses.
	reSeparators = regexp.MustCompile(`[\s_]+`)

	// reVTag strips a Bedrock inference-profile version tag, e.g. the "-v2"
	// in "claude-3-5-sonnet-20241022-v2".
	reVTag = regexp.MustCompile(`(?i)-v\d+$`)

	// reAnthropicPrefix strips a Bedrock region + provider prefix, e.g.
	// "us.anthropic." or "anthropic.", however it's cased.
	reAnthropicPrefix = regexp.MustCompile(`(?i)anthropic\.`)

	// re1M strips a "[1m]" 1M-context marker, however it's cased.
	re1M = regexp.MustCompile(`(?i)\[1m\]`)

	// reDateToken matches an 8-digit snapshot date stamp, e.g. "20241022".
	reDateToken = regexp.MustCompile(`^\d{8}$`)

	// reVersionToken matches a 1-2 digit version component, e.g. "4" or "8".
	reVersionToken = regexp.MustCompile(`^\d{1,2}$`)

	// reAllDigits matches a numeric token that's neither a date stamp nor a
	// version component (e.g. a stray build number) — dropped rather than
	// guessed at.
	reAllDigits = regexp.MustCompile(`^\d+$`)
)

// fillerWords are tokens that describe a release channel or context variant
// rather than a family or version — dropped so they never get mistaken for
// an unrecognized family name (see the otherFamily fallback below).
var fillerWords = map[string]bool{
	"latest":       true,
	"stable":       true,
	"preview":      true,
	"beta":         true,
	"alpha":        true,
	"experimental": true,
	"1m":           true,
}

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
	family, versionParts, ok := decode(id)
	if !ok {
		return "", false
	}
	label = family
	if len(versionParts) > 0 {
		label += " " + strings.Join(versionParts, ".")
	}
	return label, true
}

// Family returns just the canonical family name decoded from id (e.g.
// "Opus", "Sonnet", or the "Claude" fallback for a recognizably-Claude id
// with no known family word), without the version. It reports ok=false
// under the same conditions as Decode.
func Family(id string) (family string, ok bool) {
	family, _, ok = decode(id)
	return family, ok
}

// decode is the shared implementation behind Decode and Family: it unwraps
// known provider framing and tokenizes id into a family word and ordered
// version components.
func decode(id string) (family string, versionParts []string, ok bool) {
	s := strings.TrimSpace(id)
	if s == "" {
		return "", nil, false
	}

	// Some gateways (and env-derived IDs) use underscores or literal
	// whitespace instead of dashes.
	s = reSeparators.ReplaceAllString(s, "-")

	// Strip an OpenRouter-style "provider/" prefix, or a full Vertex
	// resource path ("projects/.../publishers/anthropic/models/claude-...").
	if i := strings.LastIndex(s, "/"); i >= 0 {
		s = s[i+1:]
	}

	// Strip a Bedrock region + provider prefix: "us.anthropic." or
	// "anthropic." both collapse to nothing. This also covers a full
	// Bedrock ARN ("arn:aws:bedrock:...:foundation-model/anthropic.claude-...")
	// once the slash-prefix strip above has dropped everything before it.
	if loc := reAnthropicPrefix.FindStringIndex(s); loc != nil {
		s = s[loc[1]:]
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

	var otherFamily string
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
		case fillerWords[lower]:
			// Release-channel adjective or context marker — not a family.
		case reAllDigits.MatchString(tok):
			// Some other numeric tag (e.g. a build number) — not a version
			// or a family; safer to drop than to guess.
		case otherFamily == "":
			otherFamily = strings.ToUpper(lower[:1]) + lower[1:]
		}
	}

	if family == "" {
		if !sawClaude {
			return "", nil, false
		}
		if otherFamily != "" {
			family = otherFamily
		} else {
			family = "Claude"
		}
	}

	// A trailing zero version component is a legacy alias artifact
	// ("claude-opus-4-0" is "Claude Opus 4", not "Claude Opus 4.0").
	for len(versionParts) > 1 && versionParts[len(versionParts)-1] == "0" {
		versionParts = versionParts[:len(versionParts)-1]
	}

	return family, versionParts, true
}
