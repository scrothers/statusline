package segment

import (
	"strings"
	"testing"

	"github.com/scrothers/statusline/internal/input"
)

func TestModelSegment(t *testing.T) {
	t.Parallel()

	t.Run("renders display name", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Model: &input.Model{DisplayName: "Opus"}}, nil)
		chunks, ok := modelSegment{}.Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "Opus") {
			t.Errorf("rendered text = %q, want it to contain Opus", chunkText(chunks))
		}
	})

	t.Run("absent model is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, nil)
		_, ok := modelSegment{}.Render(rc)
		if ok {
			t.Error("Render() ok = true, want false for nil Model")
		}
	})

	t.Run("decodes a reordered gateway id over a garbled display name", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Model: &input.Model{
			ID:          "claude-4-8-opus[1m]",
			DisplayName: "claude-4-8-opus[1m]",
		}}, nil)
		chunks, ok := modelSegment{}.Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "Opus 4.8") {
			t.Errorf("rendered text = %q, want it to contain %q", chunkText(chunks), "Opus 4.8")
		}
	})

	t.Run("icon color is per-family, not one flat accent", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, nil)

		opusChunks, ok := modelSegment{}.Render(&RenderContext{
			Payload: &input.Payload{Model: &input.Model{ID: "claude-opus-4-8"}},
			Config:  rc.Config, Theme: rc.Theme, Columns: rc.Columns,
		})
		if !ok {
			t.Fatal("Render() for opus ok = false, want true")
		}
		sonnetChunks, ok := modelSegment{}.Render(&RenderContext{
			Payload: &input.Payload{Model: &input.Model{ID: "claude-sonnet-4-6"}},
			Config:  rc.Config, Theme: rc.Theme, Columns: rc.Columns,
		})
		if !ok {
			t.Fatal("Render() for sonnet ok = false, want true")
		}

		opusIconFG, sonnetIconFG := opusChunks[0].FG, sonnetChunks[0].FG
		if opusIconFG != rc.Theme.IdentityColorFor("Opus") {
			t.Errorf("opus icon FG = %+v, want IdentityColorFor(Opus) = %+v", opusIconFG, rc.Theme.IdentityColorFor("Opus"))
		}
		if sonnetIconFG != rc.Theme.IdentityColorFor("Sonnet") {
			t.Errorf("sonnet icon FG = %+v, want IdentityColorFor(Sonnet) = %+v", sonnetIconFG, rc.Theme.IdentityColorFor("Sonnet"))
		}
		if opusIconFG == sonnetIconFG {
			t.Error("opus and sonnet icon colors are identical, want distinct per-family accents")
		}

		// Label text color is unaffected by family — stays the theme's
		// neutral identity text color for every family.
		if opusChunks[1].FG != rc.Theme.IdentityText || sonnetChunks[1].FG != rc.Theme.IdentityText {
			t.Error("label text FG changed with family, want it to stay IdentityText")
		}
	})

	t.Run("unrecognized id falls back to the flat identity accent", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Model: &input.Model{DisplayName: "Something"}}, nil)
		chunks, ok := modelSegment{}.Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if chunks[0].FG != rc.Theme.IdentityAccent {
			t.Errorf("icon FG = %+v, want the flat IdentityAccent %+v", chunks[0].FG, rc.Theme.IdentityAccent)
		}
	})
}
