package segment

import (
	"strings"
	"testing"

	"github.com/scrothers/statusline/internal/input"
	"github.com/scrothers/statusline/internal/style"
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

	t.Run("number uses a neutral identity color, not the review-state color", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{PR: &input.PR{Number: 42, ReviewState: "approved"}}, nil)
		chunks, ok := (prSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if len(chunks) < 2 {
			t.Fatal("Render() produced fewer than 2 chunks")
		}
		if chunks[1].FG != rc.Theme.IdentityText {
			t.Errorf("number FG = %+v, want theme.IdentityText %+v", chunks[1].FG, rc.Theme.IdentityText)
		}
		if chunks[1].FG == rc.Theme.Success {
			t.Errorf("number FG = %+v, should not equal the review-state color", chunks[1].FG)
		}
	})

	t.Run("review-state word still shares the review-state color", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{PR: &input.PR{Number: 42, ReviewState: "approved"}}, nil)
		chunks, ok := (prSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if len(chunks) != 3 {
			t.Fatalf("Render() produced %d chunks, want 3 (icon, number, word)", len(chunks))
		}
		if chunks[2].FG != rc.Theme.Success {
			t.Errorf("word FG = %+v, want theme.Success %+v", chunks[2].FG, rc.Theme.Success)
		}
	})

	t.Run("a URL wraps the number in a hyperlink", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{
			PR: &input.PR{Number: 42, URL: "https://github.com/scrothers/statusline/pull/42"},
		}, nil)
		chunks, ok := (prSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "https://github.com/scrothers/statusline/pull/42") {
			t.Errorf("rendered text = %q, want it to contain the PR URL in a hyperlink", chunkText(chunks))
		}
	})
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

	modeTests := []struct {
		mode        string
		wantColorOf func(rc *RenderContext) style.Color
	}{
		{mode: "INSERT", wantColorOf: func(rc *RenderContext) style.Color { return rc.Theme.Success }},
		{mode: "VISUAL", wantColorOf: func(rc *RenderContext) style.Color { return rc.Theme.Warning }},
		{mode: "VISUAL LINE", wantColorOf: func(rc *RenderContext) style.Color { return rc.Theme.Warning }},
		{mode: "NORMAL", wantColorOf: func(rc *RenderContext) style.Color { return rc.Theme.Info }},
	}
	for _, tt := range modeTests {
		t.Run("mode "+tt.mode, func(t *testing.T) {
			t.Parallel()
			rc := newTestContext(t, &input.Payload{Vim: &input.Vim{Mode: tt.mode}}, nil)
			chunks, ok := (vimSegment{}).Render(rc)
			if !ok {
				t.Fatal("Render() ok = false, want true")
			}
			if chunks[0].FG != tt.wantColorOf(rc) {
				t.Errorf("FG = %+v, want %+v", chunks[0].FG, tt.wantColorOf(rc))
			}
		})
	}
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
