package render

import (
	"strings"
	"testing"

	"github.com/scrothers/statusline/internal/config"
	"github.com/scrothers/statusline/internal/gitstatus"
	"github.com/scrothers/statusline/internal/input"
	"github.com/scrothers/statusline/internal/segment"
	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

func testRenderContext(t *testing.T, payload *input.Payload, columns int) *segment.RenderContext {
	t.Helper()
	registry, err := theme.LoadRegistry()
	if err != nil {
		t.Fatalf("theme.LoadRegistry() error = %v", err)
	}
	th, _ := theme.Resolve(registry, theme.DefaultName)
	cfg := config.Default()
	return &segment.RenderContext{
		Payload: payload,
		Config:  &cfg,
		Theme:   &th,
		Columns: columns,
	}
}

func fullPayload() *input.Payload {
	return &input.Payload{
		Model:     &input.Model{DisplayName: "Opus"},
		Workspace: &input.Workspace{CurrentDir: "/home/user/code/statusline"},
		Cost:      &input.Cost{TotalCostUSD: 1.23, TotalDurationMS: 754_000},
		ContextWindow: &input.ContextWindow{
			UsedPercentage: new(float64(42)),
		},
		RateLimits: &input.RateLimits{
			FiveHour: &input.RateLimitWindow{UsedPercentage: 30},
			SevenDay: &input.RateLimitWindow{UsedPercentage: 71},
		},
		PR:     &input.PR{Number: 128, ReviewState: "approved"},
		Vim:    &input.Vim{Mode: "INSERT"},
		Agent:  &input.Agent{Name: "reviewer"},
		Effort: &input.Effort{Level: "high"},
	}
}

func TestRender_fullPayload(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	rc := testRenderContext(t, fullPayload(), 200)
	rc.Git = &gitstatus.Status{Branch: "main", Staged: 2, Modified: 1}

	got := Render(rc, segment.Registry())

	for _, want := range []string{"Opus", "statusline", "main", "128", "42%", "$1.23", "30%", "71%", "INSERT", "reviewer", "high"} {
		if !strings.Contains(got, want) {
			t.Errorf("Render() missing %q in:\n%s", want, got)
		}
	}

	lines := strings.Split(got, "\n")
	if len(lines) != 3 {
		t.Errorf("Render() produced %d lines, want 3 (identity/repo/vitals): %q", len(lines), got)
	}
}

func TestRender_minimalPayload(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	rc := testRenderContext(t, &input.Payload{Model: &input.Model{DisplayName: "Sonnet"}, CWD: "/tmp/scratch"}, 200)
	rc.Git = &gitstatus.Status{NotARepo: true}

	got := Render(rc, segment.Registry())

	if !strings.Contains(got, "Sonnet") || !strings.Contains(got, "scratch") {
		t.Errorf("Render() = %q, missing model/directory", got)
	}
	if strings.Contains(got, "$") {
		t.Errorf("Render() = %q, cost should be absent with no Cost data", got)
	}
	// No git repo and no cost/context-window/rate-limit data: lines 2 and 3
	// both have nothing to show and disappear entirely rather than rendering
	// empty, leaving just line 1 (model + directory).
	lines := strings.Split(got, "\n")
	if len(lines) != 1 {
		t.Errorf("Render() produced %d lines, want 1 (only identity line): %q", len(lines), got)
	}
}

func TestRender_disabledLineOmitted(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	rc := testRenderContext(t, fullPayload(), 200)
	rc.Git = &gitstatus.Status{Branch: "main"}
	rc.Config.Lines[2].Enabled = false // disable the vitals line

	got := Render(rc, segment.Registry())
	if strings.Contains(got, "42%") {
		t.Errorf("Render() = %q, disabled line 3 should not render", got)
	}
	if !strings.Contains(got, "Opus") || !strings.Contains(got, "main") {
		t.Errorf("Render() = %q, lines 1/2 should still render", got)
	}
}

func TestRender_disabledSegmentOmitted(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	rc := testRenderContext(t, fullPayload(), 200)
	rc.Git = &gitstatus.Status{Branch: "main"}
	disabled := false
	rc.Config.Segments = map[string]config.SegmentConfig{"pr": {Enabled: &disabled}}

	got := Render(rc, segment.Registry())
	if strings.Contains(got, "128") {
		t.Errorf("Render() = %q, disabled pr segment should not render", got)
	}
}

func TestRender_narrowTerminalDropsBadgesButKeepsIdentity(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	rc := testRenderContext(t, fullPayload(), 20) // very narrow
	rc.Git = &gitstatus.Status{Branch: "main"}

	got := Render(rc, segment.Registry())
	if !strings.Contains(got, "Opus") {
		t.Errorf("Render() = %q, model must survive even at 20 columns", got)
	}
	if strings.Contains(got, "reviewer") {
		t.Errorf("Render() = %q, bonus badge (agent) should be dropped at 20 columns", got)
	}
}

func TestRender_noLinesEnabledProducesEmptyString(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	rc := testRenderContext(t, fullPayload(), 200)
	for i := range rc.Config.Lines {
		rc.Config.Lines[i].Enabled = false
	}

	if got := Render(rc, segment.Registry()); got != "" {
		t.Errorf("Render() = %q, want empty string when no lines are enabled", got)
	}
}

func TestSafeRender_recoversPanickingSegment(t *testing.T) {
	t.Parallel()
	rc := &segment.RenderContext{Payload: &input.Payload{}}
	chunks := safeRender(panickingSegment{}, rc)
	if chunks != nil {
		t.Errorf("safeRender() = %v, want nil after a panic", chunks)
	}
}

type panickingSegment struct{}

func (panickingSegment) ID() string    { return "panicking" }
func (panickingSegment) Priority() int { return 0 }
func (panickingSegment) Render(*segment.RenderContext) ([]style.Chunk, bool) {
	panic("boom")
}

// fakeSegment renders a single fixed chunk, for tests that need to control
// exactly which segments share a background color.
type fakeSegment struct {
	id       string
	priority int
	text     string
	bg       style.Color
}

func (f fakeSegment) ID() string    { return f.id }
func (f fakeSegment) Priority() int { return f.priority }
func (f fakeSegment) Render(*segment.RenderContext) ([]style.Chunk, bool) {
	return []style.Chunk{{Text: f.text, BG: f.bg}}, true
}

func TestRenderLine_gapMarkerForcesBreakBetweenSameBGSegments(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	bg := style.RGB(9, 9, 9)
	registry := map[string]segment.Segment{
		"alpha": fakeSegment{id: "alpha", priority: 90, text: "ALPHA", bg: bg},
		"beta":  fakeSegment{id: "beta", priority: 80, text: "BETA", bg: bg},
	}
	cfg := config.Config{Separator: config.SeparatorConfig{Style: "rounded"}}
	th := theme.Theme{}
	rc := &segment.RenderContext{Payload: &input.Payload{}, Config: &cfg, Theme: &th, Columns: 200}

	withoutGap := renderLine(rc, registry, []string{"alpha", "beta"})
	withGap := renderLine(rc, registry, []string{"alpha", GapMarker, "beta"})

	if strings.Contains(withoutGap, " ") {
		t.Errorf("chained same-bg segments should have no plain space: %q", withoutGap)
	}
	if !strings.Contains(withGap, " ") {
		t.Errorf("%q marker should introduce a plain space: %q", GapMarker, withGap)
	}
	for _, got := range []string{withoutGap, withGap} {
		if !strings.Contains(got, "ALPHA") || !strings.Contains(got, "BETA") {
			t.Errorf("renderLine() missing content: %q", got)
		}
	}
}

func TestRenderLine_gapMarkerAloneProducesNothing(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	rc := &segment.RenderContext{Payload: &input.Payload{}, Config: &config.Config{}, Theme: &theme.Theme{}, Columns: 200}
	if got := renderLine(rc, segment.Registry(), []string{GapMarker}); got != "" {
		t.Errorf("renderLine() with only a gap marker = %q, want empty", got)
	}
}
