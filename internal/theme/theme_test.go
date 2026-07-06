package theme

import (
	"testing"

	"github.com/scrothers/statusline/internal/style"
)

func TestLoadRegistry(t *testing.T) {
	t.Parallel()

	registry, err := LoadRegistry()
	if err != nil {
		t.Fatalf("LoadRegistry() error = %v", err)
	}

	want := Names()
	if len(registry) != len(want) {
		t.Fatalf("LoadRegistry() returned %d themes, want %d", len(registry), len(want))
	}
	for _, name := range want {
		th, ok := registry[name]
		if !ok {
			t.Errorf("registry missing theme %q", name)
			continue
		}
		if th.Name != name {
			t.Errorf("registry[%q].Name = %q, want %q", name, th.Name, name)
		}
		if !th.Success.Valid || !th.Warning.Valid || !th.Danger.Valid || !th.Info.Valid || !th.Muted.Valid {
			t.Errorf("registry[%q] has an unparsed semantic color: %+v", name, th)
		}
	}
}

func TestNames(t *testing.T) {
	t.Parallel()

	got := Names()
	if len(got) != 7 {
		t.Fatalf("Names() = %v, want 7 entries", got)
	}
	if got[0] != DefaultName {
		t.Errorf("Names()[0] = %q, want DefaultName %q (default listed first)", got[0], DefaultName)
	}

	// Mutating the returned slice must not affect the next call's result.
	got[0] = "mutated"
	if again := Names(); again[0] != DefaultName {
		t.Errorf("Names() after mutating a prior result = %v, want unaffected", again)
	}
}

func TestResolve(t *testing.T) {
	t.Parallel()

	registry, err := LoadRegistry()
	if err != nil {
		t.Fatalf("LoadRegistry() error = %v", err)
	}

	tests := []struct {
		name        string
		in          string
		wantName    string
		wantWarning bool
	}{
		{name: "empty defaults to gruvbox", in: "", wantName: DefaultName},
		{name: "known theme", in: "nord", wantName: "nord"},
		{name: "unknown theme falls back with warning", in: "not-a-theme", wantName: DefaultName, wantWarning: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			th, warning := Resolve(registry, tt.in)
			if th.Name != tt.wantName {
				t.Errorf("Resolve(%q) name = %q, want %q", tt.in, th.Name, tt.wantName)
			}
			if (warning != "") != tt.wantWarning {
				t.Errorf("Resolve(%q) warning = %q, wantWarning %v", tt.in, warning, tt.wantWarning)
			}
		})
	}
}

func TestThemeWithOverrides(t *testing.T) {
	t.Parallel()

	registry, err := LoadRegistry()
	if err != nil {
		t.Fatalf("LoadRegistry() error = %v", err)
	}
	base := registry[DefaultName]

	t.Run("overrides a recognized token without touching others", func(t *testing.T) {
		t.Parallel()
		got, err := base.WithOverrides(map[string]string{"success": "#00ff00"})
		if err != nil {
			t.Fatalf("WithOverrides() error = %v", err)
		}
		want, err := style.ParseHex("#00ff00")
		if err != nil {
			t.Fatalf("style.ParseHex() error = %v", err)
		}
		if got.Success != want {
			t.Errorf("Success = %+v, want %+v", got.Success, want)
		}
		if got.Warning != base.Warning {
			t.Errorf("Warning changed unexpectedly: got %+v, base %+v", got.Warning, base.Warning)
		}
		if base.Success == got.Success {
			t.Errorf("base theme was mutated by WithOverrides")
		}
	})

	t.Run("unknown token name errors", func(t *testing.T) {
		t.Parallel()
		if _, err := base.WithOverrides(map[string]string{"not_a_token": "#00ff00"}); err == nil {
			t.Error("WithOverrides() with unknown token = nil error, want error")
		}
	})

	t.Run("malformed hex errors", func(t *testing.T) {
		t.Parallel()
		if _, err := base.WithOverrides(map[string]string{"success": "not-hex"}); err == nil {
			t.Error("WithOverrides() with malformed hex = nil error, want error")
		}
	})

	t.Run("empty overrides is a no-op", func(t *testing.T) {
		t.Parallel()
		got, err := base.WithOverrides(nil)
		if err != nil {
			t.Fatalf("WithOverrides(nil) error = %v", err)
		}
		if got != base {
			t.Errorf("WithOverrides(nil) = %+v, want unchanged %+v", got, base)
		}
	})
}
