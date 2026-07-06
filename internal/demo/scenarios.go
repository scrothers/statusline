package demo

import (
	"time"

	"github.com/scrothers/statusline/internal/gitstatus"
	"github.com/scrothers/statusline/internal/input"
)

// defaultDemoColumns is generous enough that no segment gets dropped by
// width-pressure truncation, so a demo shows everything a scenario has to
// offer unless the caller deliberately asks for a narrow one.
const defaultDemoColumns = 200

// Scenario is a canned payload (and pre-collected git status) the demo
// command renders directly — no stdin JSON or real git repository needed.
type Scenario struct {
	Name    string
	Payload *input.Payload
	Git     *gitstatus.Status
	Columns int
}

// Scenarios returns the built-in demo scenarios, in display order.
func Scenarios() []Scenario {
	return []Scenario{
		minimalScenario(),
		fullScenario(),
		narrowScenario(),
	}
}

// Names returns the scenario names in the same order as Scenarios.
func Names() []string {
	scenarios := Scenarios()
	names := make([]string, len(scenarios))
	for i, s := range scenarios {
		names[i] = s.Name
	}
	return names
}

// ByName looks up a scenario by name.
func ByName(name string) (Scenario, bool) {
	for _, s := range Scenarios() {
		if s.Name == name {
			return s, true
		}
	}
	return Scenario{}, false
}

// minimalScenario is an early session: just a model and a directory, no
// git repo, no context/cost data yet.
func minimalScenario() Scenario {
	return Scenario{
		Name: "minimal",
		Payload: &input.Payload{
			Model: &input.Model{DisplayName: "Sonnet"},
			CWD:   "/tmp/scratch",
		},
		Git:     &gitstatus.Status{NotARepo: true},
		Columns: defaultDemoColumns,
	}
}

// fullScenario exercises every default-layout segment at once: the Claude
// line (model/thinking/effort/context/rate-limits/cache), the session line
// (session name/directory/lines changed/token counts/cost/duration), and the
// git line (repo/PR/branch+status/worktree).
func fullScenario() Scenario {
	return Scenario{
		Name: "full",
		Payload: &input.Payload{
			SessionName: "big-refactor",
			Model:       &input.Model{DisplayName: "Opus"},
			Workspace: &input.Workspace{
				CurrentDir: "/home/user/code/statusline",
				Repo:       &input.Repo{Host: "github.com", Owner: "scrothers", Name: "statusline"},
			},
			Cost: &input.Cost{TotalCostUSD: 2.17, TotalDurationMS: 5_025_000, TotalLinesAdded: 342, TotalLinesRemoved: 58},
			ContextWindow: &input.ContextWindow{
				UsedPercentage:    new(float64(68)),
				ContextWindowSize: 200_000,
				TotalInputTokens:  136_000,
				TotalOutputTokens: 24_000,
				CurrentUsage: &input.Usage{
					InputTokens:              20_000,
					OutputTokens:             4_500,
					CacheCreationInputTokens: 8_000,
					CacheReadInputTokens:     108_000,
				},
			},
			RateLimits: &input.RateLimits{
				FiveHour: &input.RateLimitWindow{UsedPercentage: 42, ResetsAt: time.Now().Add(2*time.Hour + 34*time.Minute).Unix()},
				SevenDay: &input.RateLimitWindow{UsedPercentage: 71, ResetsAt: time.Now().Add(3*24*time.Hour + 2*time.Hour).Unix()},
			},
			Thinking: &input.Thinking{Enabled: true},
			Effort:   &input.Effort{Level: "high"},
			PR:       &input.PR{Number: 128, URL: "https://github.com/scrothers/statusline/pull/128", ReviewState: "approved"},
			Worktree: &input.Worktree{Name: "my-feature"},
		},
		Git:     &gitstatus.Status{Branch: "main", Staged: 2, Modified: 1, Untracked: 3, Ahead: 1},
		Columns: defaultDemoColumns,
	}
}

// narrowScenario reuses fullScenario's payload at a width that forces
// priority-based truncation, so bonus badges and eventually more drop out.
func narrowScenario() Scenario {
	s := fullScenario()
	s.Name = "narrow"
	s.Columns = 30
	return s
}
