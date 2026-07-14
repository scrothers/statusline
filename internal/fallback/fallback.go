package fallback

import (
	"path/filepath"

	"github.com/scrothers/statusline/internal/input"
	"github.com/scrothers/statusline/internal/modelid"
)

// Line returns a minimal, dependency-free statusline: just the model name
// and/or directory basename, whichever is available. It never panics, even
// when payload is nil (a stdin parse failure means there's no payload at
// all).
func Line(payload *input.Payload) string {
	if payload == nil {
		return "[statusline]"
	}

	model := ""
	if payload.Model != nil {
		model = modelid.Label(payload.Model.ID, payload.Model.DisplayName)
	}

	dir := payload.CWD
	if payload.Workspace != nil && payload.Workspace.CurrentDir != "" {
		dir = payload.Workspace.CurrentDir
	}
	if dir != "" {
		dir = filepath.Base(dir)
	}

	switch {
	case model != "" && dir != "":
		return "[" + model + "] " + dir
	case model != "":
		return "[" + model + "]"
	case dir != "":
		return dir
	default:
		return "[statusline]"
	}
}
