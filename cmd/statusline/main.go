// Command statusline is a Claude Code statusLine command: it reads session
// JSON from stdin and prints a themed, Nerd-Font-powerline status line to
// stdout. See https://code.claude.com/docs/en/statusline.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/scrothers/statusline/internal/config"
	"github.com/scrothers/statusline/internal/fallback"
	"github.com/scrothers/statusline/internal/gitstatus"
	"github.com/scrothers/statusline/internal/input"
	"github.com/scrothers/statusline/internal/render"
	"github.com/scrothers/statusline/internal/segment"
	"github.com/scrothers/statusline/internal/theme"
)

// version is overridden at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	var payload *input.Payload

	// The hook invocation path must always print something and exit 0,
	// however it fails, since a blank/crashed statusline is worse than a
	// minimal degraded one. This is the outermost safety net; safeRender
	// inside internal/render catches individual segment panics before they
	// ever reach here.
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "statusline: panic: %v\n", r)
			fmt.Println(fallback.Line(payload))
		}
	}()

	configPath := flag.String("config", "", "path to a TOML config file (overrides discovery)")
	themeName := flag.String("theme", "", "theme override: gruvbox, catppuccin-mocha, tokyo-night, nord, dracula")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println("statusline " + version)
		return
	}

	var err error
	payload, err = input.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "statusline: %v\n", err)
		fmt.Println(fallback.Line(nil))
		return
	}

	cfg, warnings := config.Load(*configPath)
	for _, w := range warnings {
		fmt.Fprintln(os.Stderr, w)
	}
	if *themeName != "" {
		cfg.Theme = *themeName
	}

	th := resolveTheme(cfg)

	ctx := context.Background()
	if cfg.Budget.TotalTimeoutMS > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(cfg.Budget.TotalTimeoutMS)*time.Millisecond)
		defer cancel()
	}

	rc := &segment.RenderContext{
		Payload: payload,
		Config:  &cfg,
		Theme:   &th,
		Columns: render.Columns(),
		Now:     time.Now(),
		Git:     collectGitStatus(ctx, cfg, payload),
	}

	output := render.Render(rc, segment.Registry())
	if output == "" {
		output = fallback.Line(payload)
	}
	fmt.Println(output)
}

// resolveTheme picks and applies overrides to the configured theme, falling
// back to theme.DefaultName on any problem — an unknown theme name, an
// unreadable embedded theme registry, or a malformed override — since a
// theming issue must never block the statusline from rendering.
func resolveTheme(cfg config.Config) theme.Theme {
	registry, err := theme.LoadRegistry()
	if err != nil {
		fmt.Fprintf(os.Stderr, "statusline: %v\n", err)
		return theme.Theme{}
	}

	th, warning := theme.Resolve(registry, cfg.Theme)
	if warning != "" {
		fmt.Fprintln(os.Stderr, warning)
	}

	if len(cfg.ThemeOverrides) == 0 {
		return th
	}
	overridden, err := th.WithOverrides(cfg.ThemeOverrides)
	if err != nil {
		fmt.Fprintf(os.Stderr, "statusline: %v\n", err)
	}
	return overridden
}

// collectGitStatus returns nil when git collection is disabled or there's
// no directory to check; a nil *gitstatus.Status makes the git segment
// (and PR badge, which piggybacks on the same repo context) simply omit
// itself rather than erroring.
func collectGitStatus(ctx context.Context, cfg config.Config, payload *input.Payload) *gitstatus.Status {
	if !cfg.Git.IsEnabled() {
		return nil
	}

	dir := payload.CWD
	if payload.Workspace != nil && payload.Workspace.CurrentDir != "" {
		dir = payload.Workspace.CurrentDir
	}
	if dir == "" {
		return nil
	}

	st, err := gitstatus.Collect(ctx, cfg.Git, payload.SessionID, dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "statusline: %v\n", err)
		return nil
	}
	return &st
}
