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
