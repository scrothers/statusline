package input

import (
	"encoding/json"
	"fmt"
	"io"
)

// maxPayloadBytes bounds how much of stdin Parse will read, guarding against
// a pathological or hostile input source.
const maxPayloadBytes = 1 << 20 // 1MB

// Parse decodes a Payload from r. Unknown fields are ignored rather than
// rejected, since Claude Code adds schema fields over time and this binary
// must keep working against payloads newer than its own schema knowledge.
func Parse(r io.Reader) (*Payload, error) {
	dec := json.NewDecoder(io.LimitReader(r, maxPayloadBytes))
	var p Payload
	if err := dec.Decode(&p); err != nil {
		return nil, fmt.Errorf("input: parse: %w", err)
	}
	return &p, nil
}
