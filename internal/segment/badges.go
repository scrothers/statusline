package segment

import (
	"fmt"

	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// prSegment renders the open pull request for the current branch, colored
// by review state. It shares line 2's chrome (a real pill), unlike the
// low-priority badges below (which render with no pill background).
type prSegment struct{}

func (prSegment) ID() string { return "pr" }

func (prSegment) Priority() int { return 70 }

func (prSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	pr := rc.Payload.PR
	if pr == nil {
		return nil, false
	}

	nerd := rc.Config.NerdFontEnabled()
	bg := rc.Theme.Line2Bg

	color, iconKey := rc.Theme.Info, theme.IconPR
	switch pr.ReviewState {
	case "draft":
		color, iconKey = rc.Theme.Muted, theme.IconPRDraft
	case "pending":
		color, iconKey = rc.Theme.Warning, theme.IconPRPending
	case "approved":
		color, iconKey = rc.Theme.Success, theme.IconPRApproved
	case "changes_requested":
		color, iconKey = rc.Theme.Danger, theme.IconPRChangesRequest
	}

	text := fmt.Sprintf(" %s #%d ", theme.Glyph(iconKey, nerd), pr.Number)
	return []style.Chunk{{Text: text, FG: color, BG: bg}}, true
}

// vimSegment renders the current vim mode as a plain, pill-less badge.
type vimSegment struct{}

func (vimSegment) ID() string { return "vim" }

func (vimSegment) Priority() int { return 30 }

func (vimSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	if rc.Payload.Vim == nil || rc.Payload.Vim.Mode == "" {
		return nil, false
	}

	color := rc.Theme.Info
	switch rc.Payload.Vim.Mode {
	case "INSERT":
		color = rc.Theme.Success
	case "VISUAL", "VISUAL LINE":
		color = rc.Theme.Warning
	}

	return []style.Chunk{{Text: rc.Payload.Vim.Mode, FG: color, BG: style.Default}}, true
}

// agentSegment renders the running subagent's name as a plain badge.
type agentSegment struct{}

func (agentSegment) ID() string { return "agent" }

func (agentSegment) Priority() int { return 30 }

func (agentSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	if rc.Payload.Agent == nil || rc.Payload.Agent.Name == "" {
		return nil, false
	}

	nerd := rc.Config.NerdFontEnabled()
	text := rc.Payload.Agent.Name
	if icon := theme.Glyph(theme.IconAgent, nerd); icon != "" {
		text = icon + " " + text
	}
	return []style.Chunk{{Text: text, FG: rc.Theme.Muted, BG: style.Default}}, true
}

// effortSegment renders the current reasoning effort level as plain text in
// the theme's identity accent, since effort is a setting (a fact), not a
// state — it never turns red/green.
type effortSegment struct{}

func (effortSegment) ID() string { return "effort" }

func (effortSegment) Priority() int { return 30 }

func (effortSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	if rc.Payload.Effort == nil || rc.Payload.Effort.Level == "" {
		return nil, false
	}
	return []style.Chunk{{Text: rc.Payload.Effort.Level, FG: rc.Theme.IdentityAccent, BG: style.Default}}, true
}

// outputStyleSegment renders the current output style name, skipping the
// common "default" style since it carries no information.
type outputStyleSegment struct{}

func (outputStyleSegment) ID() string { return "output_style" }

func (outputStyleSegment) Priority() int { return 20 }

func (outputStyleSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	if rc.Payload.OutputStyle == nil || rc.Payload.OutputStyle.Name == "" || rc.Payload.OutputStyle.Name == "default" {
		return nil, false
	}

	nerd := rc.Config.NerdFontEnabled()
	text := rc.Payload.OutputStyle.Name
	if icon := theme.Glyph(theme.IconOutputStyle, nerd); icon != "" {
		text = icon + " " + text
	}
	return []style.Chunk{{Text: text, FG: rc.Theme.Muted, BG: style.Default}}, true
}
