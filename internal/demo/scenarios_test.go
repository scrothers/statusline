package demo

import "testing"

func TestScenarios(t *testing.T) {
	t.Parallel()

	scenarios := Scenarios()
	if len(scenarios) == 0 {
		t.Fatal("Scenarios() returned none")
	}
	for _, s := range scenarios {
		if s.Name == "" {
			t.Error("scenario has an empty Name")
		}
		if s.Payload == nil {
			t.Errorf("scenario %q has a nil Payload", s.Name)
		}
		if s.Columns <= 0 {
			t.Errorf("scenario %q has Columns = %d, want > 0", s.Name, s.Columns)
		}
	}
}

func TestNames(t *testing.T) {
	t.Parallel()

	names := Names()
	scenarios := Scenarios()
	if len(names) != len(scenarios) {
		t.Fatalf("Names() has %d entries, Scenarios() has %d", len(names), len(scenarios))
	}
	for i, s := range scenarios {
		if names[i] != s.Name {
			t.Errorf("Names()[%d] = %q, want %q", i, names[i], s.Name)
		}
	}
}

func TestByName(t *testing.T) {
	t.Parallel()

	t.Run("known scenario", func(t *testing.T) {
		t.Parallel()
		s, ok := ByName("full")
		if !ok {
			t.Fatal("ByName(full) ok = false, want true")
		}
		if s.Name != "full" {
			t.Errorf("ByName(full).Name = %q, want full", s.Name)
		}
	})

	t.Run("unknown scenario", func(t *testing.T) {
		t.Parallel()
		if _, ok := ByName("not-a-scenario"); ok {
			t.Error("ByName(unknown) ok = true, want false")
		}
	})
}

func TestNarrowScenarioReusesFullPayload(t *testing.T) {
	t.Parallel()

	full, ok := ByName("full")
	if !ok {
		t.Fatal("ByName(full) not found")
	}
	narrow, ok := ByName("narrow")
	if !ok {
		t.Fatal("ByName(narrow) not found")
	}

	if narrow.Columns >= full.Columns {
		t.Errorf("narrow.Columns = %d, want < full.Columns %d", narrow.Columns, full.Columns)
	}
	if narrow.Payload.Model.DisplayName != full.Payload.Model.DisplayName {
		t.Error("narrow scenario should reuse full scenario's payload content")
	}
}
