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

// fullPayload exercises every default-layout segment at once: the Claude
// line (model/thinking/effort/context/rate-limits/cache), the session line
// (session name/directory/lines changed/cost), and the git line
// (repo/PR/branch/worktree — branch comes from rc.Git, set by callers).
func fullPayload() *input.Payload {
	return &input.Payload{
		SessionName: "big-refactor",
		Model:       &input.Model{DisplayName: "Opus"},
		Workspace: &input.Workspace{
			CurrentDir: "/home/user/code/statusline",
			Repo:       &input.Repo{Host: "github.com", Owner: "scrothers", Name: "statusline"},
		},
		Cost: &input.Cost{TotalCostUSD: 1.23, TotalDurationMS: 754_000, TotalLinesAdded: 342, TotalLinesRemoved: 58},
		ContextWindow: &input.ContextWindow{
			UsedPercentage: new(float64(42)),
			CurrentUsage:   &input.Usage{InputTokens: 1000, CacheCreationInputTokens: 1000, CacheReadInputTokens: 8000},
		},
		RateLimits: &input.RateLimits{
			FiveHour: &input.RateLimitWindow{UsedPercentage: 30},
			SevenDay: &input.RateLimitWindow{UsedPercentage: 71},
		},
		PR:       &input.PR{Number: 128, ReviewState: "approved"},
		Thinking: &input.Thinking{Enabled: true},
		Effort:   &input.Effort{Level: "high"},
		Worktree: &input.Worktree{Name: "my-feature"},
	}
}

func TestRender_fullPayload(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	rc := testRenderContext(t, fullPayload(), 200)
	rc.Git = &gitstatus.Status{Branch: "main", Staged: 2, Modified: 1}

	got := Render(rc, segment.Registry())

	for _, want := range []string{
		"Opus", "high", "42%", "30%", "71%", "80%", "8.0k", // Claude line
		"big-refactor", "statusline", "+342", "-58", "1.23", // session line
		"github.com", "scrothers", "128", "approved", "main", "my-feature", // git line
	} {
		if !strings.Contains(got, want) {
			t.Errorf("Render() missing %q in:\n%s", want, got)
		}
	}

	lines := strings.Split(got, "\n")
	if len(lines) != 3 {
		t.Errorf("Render() produced %d lines, want 3 (Claude/session/git): %q", len(lines), got)
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
	// No git repo and nothing for the git line to show: it disappears
	// entirely rather than rendering empty, leaving the Claude line (model
	// only) and the session line (directory only).
	lines := strings.Split(got, "\n")
	if len(lines) != 2 {
		t.Errorf("Render() produced %d lines, want 2 (Claude + session, no git line): %q", len(lines), got)
	}
}

func TestRender_disabledLineOmitted(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	rc := testRenderContext(t, fullPayload(), 200)
	rc.Git = &gitstatus.Status{Branch: "main"}
	rc.Config.Lines[2].Enabled = false // disable the git information line

	got := Render(rc, segment.Registry())
	if strings.Contains(got, "main") || strings.Contains(got, "128") {
		t.Errorf("Render() = %q, disabled git line should not render", got)
	}
	if !strings.Contains(got, "Opus") || !strings.Contains(got, "big-refactor") {
		t.Errorf("Render() = %q, Claude/session lines should still render", got)
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

func TestRender_narrowTerminalDropsLowPriorityButKeepsIdentity(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	rc := testRenderContext(t, fullPayload(), 20) // very narrow
	rc.Git = &gitstatus.Status{Branch: "main"}

	got := Render(rc, segment.Registry())
	if !strings.Contains(got, "Opus") {
		t.Errorf("Render() = %q, model must survive even at 20 columns", got)
	}
	if strings.Contains(got, "8.0k") {
		t.Errorf("Render() = %q, cache (lowest priority on the Claude line) should be dropped at 20 columns", got)
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

// fakeSegment renders a single fixed chunk, for tests that don't need a
// real segment implementation.
type fakeSegment struct {
	id       string
	priority int
	text     string
}

func (f fakeSegment) ID() string    { return f.id }
func (f fakeSegment) Priority() int { return f.priority }
func (f fakeSegment) Render(*segment.RenderContext) ([]style.Chunk, bool) {
	return []style.Chunk{{Text: f.text}}, true
}

func TestRenderLine_joinsSegmentsWithADivider(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	registry := map[string]segment.Segment{
		"alpha": fakeSegment{id: "alpha", priority: 90, text: "ALPHA"},
		"beta":  fakeSegment{id: "beta", priority: 80, text: "BETA"},
	}
	cfg := config.Config{}
	th := theme.Theme{}
	rc := &segment.RenderContext{Payload: &input.Payload{}, Config: &cfg, Theme: &th, Columns: 200}

	got := renderLine(rc, registry, []string{"alpha", "beta"})

	if !strings.Contains(got, "ALPHA") || !strings.Contains(got, "BETA") {
		t.Errorf("renderLine() missing content: %q", got)
	}
	if !strings.Contains(got, " ") {
		t.Errorf("renderLine() should join segments with a divider (plain space present): %q", got)
	}
	// No background is ever painted, under NO_COLOR or otherwise: nothing in
	// the join layer should introduce color at all here since NO_COLOR strips
	// every Paint() call to plain text.
	if strings.Contains(got, "\x1b[") {
		t.Errorf("renderLine() under NO_COLOR should contain no ANSI escapes: %q", got)
	}
}

func TestRenderLine_unknownSegmentIDsAreSkipped(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	rc := &segment.RenderContext{Payload: &input.Payload{}, Config: &config.Config{}, Theme: &theme.Theme{}, Columns: 200}
	if got := renderLine(rc, segment.Registry(), []string{"not-a-real-segment-id"}); got != "" {
		t.Errorf("renderLine() with only an unknown ID = %q, want empty", got)
	}
}
