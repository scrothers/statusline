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
}
