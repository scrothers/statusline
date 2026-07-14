//go:build e2e

// Package statusline_test end-to-end tests the built statusline binary by
// executing it as a real subprocess, exactly as Claude Code would.
package statusline_test

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binPath string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "statusline-e2e-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	binPath = filepath.Join(dir, "statusline")
	build := exec.Command("go", "build", "-o", binPath, "../../../cmd/statusline")
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		panic("build statusline: " + err.Error())
	}

	os.Exit(m.Run())
}

func run(t *testing.T, stdin string, env []string, args ...string) (stdout string, exitCode int) {
	t.Helper()
	cmd := exec.Command(binPath, args...)
	cmd.Stdin = bytes.NewBufferString(stdin)
	cmd.Env = append(os.Environ(), env...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err == nil {
		return out.String(), 0
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("run statusline: %v", err)
	}
	return out.String(), exitErr.ExitCode()
}

func TestE2E_fixturesAlwaysExitZeroAndPrint(t *testing.T) {
	fixtures := map[string]string{
		"minimal early session": `{"model":{"display_name":"Opus"},"cwd":"/tmp","session_id":"s1"}`,
		"full session": `{"model":{"display_name":"Sonnet"},"workspace":{"current_dir":"/tmp"},` +
			`"context_window":{"used_percentage":80},"cost":{"total_cost_usd":1.5,"total_duration_ms":90000},` +
			`"rate_limits":{"five_hour":{"used_percentage":50,"resets_at":1},"seven_day":{"used_percentage":10,"resets_at":2}},` +
			`"vim":{"mode":"NORMAL"},"effort":{"level":"high"},"pr":{"number":9,"review_state":"approved"},"session_id":"s2"}`,
		"malformed json": `not json at all`,
		"empty stdin":    ``,
	}
	for name, fixture := range fixtures {
		t.Run(name, func(t *testing.T) {
			out, code := run(t, fixture, nil)
			if code != 0 {
				t.Errorf("exit code = %d, want 0", code)
			}
			if strings.TrimSpace(out) == "" {
				t.Error("stdout empty, want a non-empty statusline")
			}
		})
	}
}

func TestE2E_gatewayModelIDDecodesOverDisplayName(t *testing.T) {
	payload := `{"model":{"id":"claude-4-8-opus[1m]","display_name":"claude-4-8-opus[1m]"},` +
		`"cwd":"/tmp","session_id":"s6"}`
	out, code := run(t, payload, nil)
	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	if !strings.Contains(out, "Opus 4.8") {
		t.Errorf("output = %q, want it to contain %q", out, "Opus 4.8")
	}
}

func TestE2E_providerBadgeFallbackText(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(configPath, []byte("nerd_font = false\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	payload := `{"model":{"id":"anthropic.claude-opus-4-8","display_name":"anthropic.claude-opus-4-8"},` +
		`"cwd":"/tmp","session_id":"s7"}`
	out, code := run(t, payload, nil, "--config", configPath)
	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	if !strings.Contains(out, "AWS") {
		t.Errorf("output = %q, want it to contain the AWS provider badge fallback", out)
	}
}

func TestE2E_version(t *testing.T) {
	out, code := run(t, "", nil, "--version")
	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	if !strings.Contains(out, "statusline") {
		t.Errorf("output = %q, want it to contain statusline", out)
	}
}

func TestE2E_allThemesRenderCleanly(t *testing.T) {
	payload := `{"model":{"display_name":"Opus"},"cwd":"/tmp","context_window":{"used_percentage":50},"session_id":"s3"}`
	for _, name := range []string{"gruvbox", "catppuccin-mocha", "tokyo-night", "nord", "dracula", "claude-dark", "claude-light"} {
		t.Run(name, func(t *testing.T) {
			out, code := run(t, payload, nil, "--theme", name)
			if code != 0 {
				t.Errorf("exit code = %d, want 0", code)
			}
			if strings.TrimSpace(out) == "" {
				t.Error("stdout empty")
			}
		})
	}
}

func TestE2E_unknownThemeFallsBackWithoutFailing(t *testing.T) {
	payload := `{"model":{"display_name":"Opus"},"cwd":"/tmp","session_id":"s4"}`
	out, code := run(t, payload, nil, "--theme", "not-a-real-theme")
	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	if strings.TrimSpace(out) == "" {
		t.Error("stdout empty")
	}
}

func TestE2E_noColorProducesPlainText(t *testing.T) {
	payload := `{"model":{"display_name":"Opus"},"cwd":"/tmp","context_window":{"used_percentage":50},"session_id":"s5"}`
	out, code := run(t, payload, []string{"NO_COLOR=1"})
	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	if strings.Contains(out, "\x1b[") {
		t.Errorf("output contains ANSI escapes under NO_COLOR: %q", out)
	}
	if !strings.Contains(out, "Opus") {
		t.Errorf("output = %q, want it to contain Opus", out)
	}
}

func TestE2E_demoDefaultShowsAllThemes(t *testing.T) {
	out, code := run(t, "", []string{"NO_COLOR=1"}, "demo")
	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	for _, theme := range []string{"gruvbox", "catppuccin-mocha", "tokyo-night", "nord", "dracula", "claude-dark", "claude-light"} {
		if !strings.Contains(out, theme) {
			t.Errorf("demo output missing theme header %q:\n%s", theme, out)
		}
	}
	if !strings.Contains(out, "Opus") {
		t.Errorf("demo output (full scenario) missing model name:\n%s", out)
	}
}

func TestE2E_demoSingleThemeAndScenario(t *testing.T) {
	out, code := run(t, "", []string{"NO_COLOR=1"}, "demo", "--theme", "dracula", "--scenario", "minimal")
	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	if strings.Contains(out, "──") {
		t.Errorf("single-theme demo should not print a theme header:\n%s", out)
	}
	if !strings.Contains(out, "Sonnet") {
		t.Errorf("demo output (minimal scenario) missing model name:\n%s", out)
	}
	if strings.Contains(out, "$") {
		t.Errorf("minimal scenario should have no cost segment:\n%s", out)
	}
}

func TestE2E_demoUnknownScenarioErrors(t *testing.T) {
	_, code := run(t, "", nil, "demo", "--scenario", "not-a-scenario")
	if code == 0 {
		t.Error("exit code = 0, want non-zero for an unknown scenario")
	}
}

func TestE2E_demoNarrowScenarioDropsLowPriority(t *testing.T) {
	out, code := run(t, "", []string{"NO_COLOR=1"}, "demo", "--theme", "gruvbox", "--scenario", "narrow")
	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	if !strings.Contains(out, "Opus") {
		t.Errorf("model must survive even in the narrow scenario:\n%s", out)
	}
	if strings.Contains(out, "16.0k") {
		t.Errorf("cache (lowest priority on the Claude line) should be dropped in the narrow scenario:\n%s", out)
	}
}
