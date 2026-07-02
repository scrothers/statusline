package fallback

import (
	"testing"

	"github.com/scrothers/statusline/internal/input"
)

func TestLine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		payload *input.Payload
		want    string
	}{
		{name: "nil payload", payload: nil, want: "[statusline]"},
		{name: "empty payload", payload: &input.Payload{}, want: "[statusline]"},
		{
			name:    "model and cwd",
			payload: &input.Payload{Model: &input.Model{DisplayName: "Opus"}, CWD: "/home/user/project"},
			want:    "[Opus] project",
		},
		{
			name:    "model and workspace current dir preferred over cwd",
			payload: &input.Payload{Model: &input.Model{DisplayName: "Opus"}, CWD: "/wrong", Workspace: &input.Workspace{CurrentDir: "/home/user/project"}},
			want:    "[Opus] project",
		},
		{
			name:    "model only",
			payload: &input.Payload{Model: &input.Model{DisplayName: "Opus"}},
			want:    "[Opus]",
		},
		{
			name:    "directory only",
			payload: &input.Payload{CWD: "/home/user/project"},
			want:    "project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := Line(tt.payload); got != tt.want {
				t.Errorf("Line() = %q, want %q", got, tt.want)
			}
		})
	}
}
