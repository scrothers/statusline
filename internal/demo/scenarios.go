package demo

import (
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

// fullScenario exercises every segment at once: dirty git repo, an
// approved PR, context/cost/rate-limit data, and every bonus badge.
func fullScenario() Scenario {
	return Scenario{
		Name: "full",
		Payload: &input.Payload{
			Model:     &input.Model{DisplayName: "Opus"},
			Workspace: &input.Workspace{CurrentDir: "/home/user/code/statusline"},
			Cost:      &input.Cost{TotalCostUSD: 2.17, TotalDurationMS: 5_025_000},
			ContextWindow: &input.ContextWindow{
				UsedPercentage: new(float64(68)),
			},
			RateLimits: &input.RateLimits{
				FiveHour: &input.RateLimitWindow{UsedPercentage: 42},
				SevenDay: &input.RateLimitWindow{UsedPercentage: 71},
			},
			Vim:    &input.Vim{Mode: "INSERT"},
			Agent:  &input.Agent{Name: "reviewer"},
			Effort: &input.Effort{Level: "high"},
			PR:     &input.PR{Number: 128, ReviewState: "approved"},
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
