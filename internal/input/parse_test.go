package input

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		json    string
		wantErr bool
		check   func(t *testing.T, p *Payload)
	}{
		{
			name: "minimal early session",
			json: `{"cwd":"/home/user/project","session_id":"abc123","transcript_path":"/tmp/t.jsonl","version":"2.1.90","model":{"id":"claude-opus-4-8","display_name":"Opus"}}`,
			check: func(t *testing.T, p *Payload) {
				t.Helper()
				if p.CWD != "/home/user/project" {
					t.Errorf("CWD = %q, want /home/user/project", p.CWD)
				}
				if p.Model == nil || p.Model.DisplayName != "Opus" {
					t.Errorf("Model = %+v, want DisplayName=Opus", p.Model)
				}
				if p.RateLimits != nil {
					t.Errorf("RateLimits = %+v, want nil", p.RateLimits)
				}
				if p.ContextWindow != nil {
					t.Errorf("ContextWindow = %+v, want nil", p.ContextWindow)
				}
			},
		},
		{
			name: "nullable context window fields",
			json: `{"context_window":{"total_input_tokens":0,"total_output_tokens":0,"context_window_size":200000,"used_percentage":null,"remaining_percentage":null,"current_usage":null}}`,
			check: func(t *testing.T, p *Payload) {
				t.Helper()
				if p.ContextWindow == nil {
					t.Fatal("ContextWindow = nil, want non-nil")
				}
				if p.ContextWindow.UsedPercentage != nil {
					t.Errorf("UsedPercentage = %v, want nil", p.ContextWindow.UsedPercentage)
				}
				if p.ContextWindow.CurrentUsage != nil {
					t.Errorf("CurrentUsage = %v, want nil", p.ContextWindow.CurrentUsage)
				}
			},
		},
		{
			name: "unknown fields are ignored",
			json: `{"cwd":"/x","totally_new_field":{"nested":true},"model":{"id":"m","display_name":"M"}}`,
			check: func(t *testing.T, p *Payload) {
				t.Helper()
				if p.CWD != "/x" {
					t.Errorf("CWD = %q, want /x", p.CWD)
				}
			},
		},
		{
			name:    "malformed json errors",
			json:    `{not valid json`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p, err := Parse(strings.NewReader(tt.json))
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			tt.check(t, p)
		})
	}
}

// fullBenchmarkPayload is a realistic full-featured payload (every top-level
// field populated) for BenchmarkParse — representative of the actual worst
// case Parse sees on a live session, not a minimal fixture.
const fullBenchmarkPayload = `{
	"cwd": "/home/user/code/statusline",
	"session_id": "abc123",
	"session_name": "big-refactor",
	"prompt_id": "p1",
	"transcript_path": "/tmp/t.jsonl",
	"version": "2.1.90",
	"model": {"id": "claude-opus-4-8", "display_name": "Opus"},
	"workspace": {
		"current_dir": "/home/user/code/statusline",
		"project_dir": "/home/user/code/statusline",
		"added_dirs": ["/home/user/code/statusline/vendor"],
		"repo": {"host": "github.com", "owner": "scrothers", "name": "statusline"}
	},
	"output_style": {"name": "Explanatory"},
	"cost": {
		"total_cost_usd": 2.17,
		"total_duration_ms": 5025000,
		"total_api_duration_ms": 4800000,
		"total_lines_added": 342,
		"total_lines_removed": 58
	},
	"context_window": {
		"total_input_tokens": 136000,
		"total_output_tokens": 24000,
		"context_window_size": 200000,
		"used_percentage": 68,
		"remaining_percentage": 32,
		"current_usage": {
			"input_tokens": 20000,
			"output_tokens": 4500,
			"cache_creation_input_tokens": 8000,
			"cache_read_input_tokens": 108000
		}
	},
	"exceeds_200k_tokens": false,
	"effort": {"level": "high"},
	"thinking": {"enabled": true},
	"rate_limits": {
		"five_hour": {"used_percentage": 42, "resets_at": 1735689600},
		"seven_day": {"used_percentage": 71, "resets_at": 1735948800}
	},
	"vim": {"mode": "INSERT"},
	"agent": {"name": "reviewer"},
	"pr": {"number": 128, "url": "https://github.com/scrothers/statusline/pull/128", "review_state": "approved"},
	"worktree": {"name": "my-feature", "path": "/home/user/worktrees/my-feature", "branch": "my-feature"}
}`

// BenchmarkParse measures decoding a realistic full-featured payload, the
// one unavoidable per-invocation cost on the hot path from stdin to render.
func BenchmarkParse(b *testing.B) {
	for b.Loop() {
		if _, err := Parse(strings.NewReader(fullBenchmarkPayload)); err != nil {
			b.Fatalf("Parse() error = %v", err)
		}
	}
}

func TestDeref(t *testing.T) {
	t.Parallel()

	if got := Deref(nil, "fallback"); got != "fallback" {
		t.Errorf("Deref(nil, fallback) = %q, want fallback", got)
	}
	v := "present"
	if got := Deref(&v, "fallback"); got != "present" {
		t.Errorf("Deref(&v, fallback) = %q, want present", got)
	}
}
