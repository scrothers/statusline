package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/scrothers/statusline/internal/config"
	"github.com/scrothers/statusline/internal/demo"
	"github.com/scrothers/statusline/internal/render"
	"github.com/scrothers/statusline/internal/segment"
	"github.com/scrothers/statusline/internal/theme"
)

// newDemoCmd builds the `statusline demo` subcommand: it renders built-in
// sample payloads so a user can preview themes and layout without piping
// real Claude Code JSON on stdin or wiring up settings.json first.
func newDemoCmd() *cobra.Command {
	var themeName string
	var scenarioName string
	var columns int

	cmd := &cobra.Command{
		Use:   "demo",
		Short: "Render sample statuslines to preview themes and layout",
		Long: "demo renders built-in sample payloads — no stdin JSON or real git\n" +
			"repository needed — so you can see what a theme or layout looks like\n" +
			"before wiring statusline into Claude Code.",
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runDemo(themeName, scenarioName, columns)
		},
	}
	cmd.Flags().StringVar(&themeName, "theme", "", "theme to preview (default: all built-in themes)")
	cmd.Flags().StringVar(&scenarioName, "scenario", "full",
		fmt.Sprintf("scenario to render: %s", strings.Join(demo.Names(), ", ")))
	cmd.Flags().IntVar(&columns, "columns", 0, "override the scenario's terminal width (0 = scenario default)")
	return cmd
}

func runDemo(themeName, scenarioName string, columns int) error {
	scenario, ok := demo.ByName(scenarioName)
	if !ok {
		return fmt.Errorf("statusline: unknown scenario %q (want one of: %s)", scenarioName, strings.Join(demo.Names(), ", "))
	}
	if columns > 0 {
		scenario.Columns = columns
	}

	registry, err := theme.LoadRegistry()
	if err != nil {
		return err
	}

	names := theme.Names()
	if themeName != "" {
		names = []string{themeName}
	}

	cfg := config.Default()
	reg := segment.Registry()
	for i, name := range names {
		th, warning := theme.Resolve(registry, name)
		if warning != "" {
			fmt.Fprintln(os.Stderr, warning)
		}
		if len(names) > 1 {
			fmt.Printf("── %s (%s) ──\n", th.Name, scenario.Name)
		}

		rc := &segment.RenderContext{
			Payload: scenario.Payload,
			Config:  &cfg,
			Theme:   &th,
			Columns: scenario.Columns,
			Git:     scenario.Git,
		}
		fmt.Println(render.Render(rc, reg))
		if len(names) > 1 && i < len(names)-1 {
			fmt.Println()
		}
	}
	return nil
}
