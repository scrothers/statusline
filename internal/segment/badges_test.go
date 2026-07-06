package segment

import (
	"strings"
	"testing"

	"github.com/scrothers/statusline/internal/input"
)

func TestPRSegment(t *testing.T) {
	t.Parallel()

	t.Run("absent PR is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, nil)
		if _, ok := (prSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false for nil PR")
		}
	})

	tests := []struct {
		reviewState string
		wantColorOf func(rc *RenderContext) any
	}{
		{reviewState: "draft", wantColorOf: func(rc *RenderContext) any { return rc.Theme.Muted }},
		{reviewState: "pending", wantColorOf: func(rc *RenderContext) any { return rc.Theme.Warning }},
		{reviewState: "approved", wantColorOf: func(rc *RenderContext) any { return rc.Theme.Success }},
		{reviewState: "changes_requested", wantColorOf: func(rc *RenderContext) any { return rc.Theme.Danger }},
		{reviewState: "", wantColorOf: func(rc *RenderContext) any { return rc.Theme.Info }},
	}
	for _, tt := range tests {
		t.Run("review state "+tt.reviewState, func(t *testing.T) {
			t.Parallel()
			rc := newTestContext(t, &input.Payload{PR: &input.PR{Number: 42, ReviewState: tt.reviewState}}, nil)
			chunks, ok := (prSegment{}).Render(rc)
			if !ok {
				t.Fatal("Render() ok = false, want true")
			}
			if !strings.Contains(chunkText(chunks), "42") {
				t.Errorf("rendered text = %q, want it to contain the PR number", chunkText(chunks))
			}
			if any(chunks[0].FG) != tt.wantColorOf(rc) {
				t.Errorf("FG = %+v, want %+v", chunks[0].FG, tt.wantColorOf(rc))
			}
		})
	}
}

func TestVimSegment(t *testing.T) {
	t.Parallel()

	t.Run("absent vim is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, nil)
		if _, ok := (vimSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false for nil Vim")
		}
	})

	t.Run("renders mode text and has no pill background", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Vim: &input.Vim{Mode: "INSERT"}}, nil)
		chunks, ok := (vimSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if chunkText(chunks) != "INSERT" {
			t.Errorf("rendered text = %q, want INSERT", chunkText(chunks))
		}
		if chunks[0].BG.Valid {
			t.Errorf("BG = %+v, want the terminal default (no pill)", chunks[0].BG)
		}
	})
}

func TestAgentSegment(t *testing.T) {
	t.Parallel()

	t.Run("absent agent is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, nil)
		if _, ok := (agentSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false for nil Agent")
		}
	})

	t.Run("renders agent name", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Agent: &input.Agent{Name: "reviewer"}}, nil)
		chunks, ok := (agentSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "reviewer") {
			t.Errorf("rendered text = %q, want it to contain reviewer", chunkText(chunks))
		}
	})
}

func TestOutputStyleSegment(t *testing.T) {
	t.Parallel()

	t.Run("absent output style is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, nil)
		if _, ok := (outputStyleSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false for nil OutputStyle")
		}
	})

	t.Run("default style is omitted as uninformative", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{OutputStyle: &input.OutputStyle{Name: "default"}}, nil)
		if _, ok := (outputStyleSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false for the default style")
		}
	})

	t.Run("non-default style renders", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{OutputStyle: &input.OutputStyle{Name: "Explanatory"}}, nil)
		chunks, ok := (outputStyleSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "Explanatory") {
			t.Errorf("rendered text = %q, want it to contain Explanatory", chunkText(chunks))
		}
	})
}
