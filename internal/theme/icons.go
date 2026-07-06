package theme

// Icon segment keys. Exported so segment code references a stable
// identifier rather than a magic string.
const (
	IconModel             = "model"
	IconDirectory         = "directory"
	IconGitBranch         = "git_branch"
	IconGitStaged         = "git_staged"
	IconGitModified       = "git_modified"
	IconGitUntracked      = "git_untracked"
	IconGitAhead          = "git_ahead"
	IconGitBehind         = "git_behind"
	IconGitCleanDot       = "git_clean_dot"
	IconPR                = "pr"
	IconPRDraft           = "pr_draft"
	IconPRPending         = "pr_pending"
	IconPRApproved        = "pr_approved"
	IconPRChangesRequest  = "pr_changes_requested"
	IconContextWindow     = "context_window"
	IconContextAlert      = "context_alert"
	IconCost              = "cost"
	IconDuration          = "duration"
	IconRateLimitFiveHour = "ratelimit_5h"
	IconRateLimitWeek     = "ratelimit_7d"
	IconVim               = "vim"
	IconAgent             = "agent"
	IconOutputStyle       = "output_style"
	IconThinking          = "thinking"
	IconCache             = "cache"
	IconSessionName       = "session_name"
	IconRepo              = "repo"
	IconWorktree          = "worktree"
)

// Icon pairs a Nerd Font glyph with a plain-Unicode/ASCII fallback for
// terminals without a Nerd Font installed. Codepoints are sourced from the
// Nerd Fonts cheat sheet by icon name; verify against nerdfonts.com if a
// glyph renders as tofu, since exact codepoints have drifted across nerd
// font releases (e.g. the Material Design Icons migration to the U+F0000+
// supplementary plane) while names stay stable. Glyphs are written as
// explicit \u/\U escapes rather than literal source bytes so the private-use
// codepoints survive editing/transport intact.
type Icon struct {
	Glyph    string
	Fallback string
}

// Icons is the shared, theme-independent icon table: only color varies per
// theme, not iconography.
var Icons = map[string]Icon{
	IconModel:             {Glyph: "\U000F06A9", Fallback: "AI"},        // nf-md-robot_outline
	IconDirectory:         {Glyph: "\U000F024B", Fallback: "~"},         // nf-md-folder
	IconGitBranch:         {Glyph: "\uE725", Fallback: "git:"},          // nf-dev-git_branch
	IconGitStaged:         {Glyph: "\U000F0521", Fallback: "+"},         // nf-md-plus_circle_outline
	IconGitModified:       {Glyph: "\uF040", Fallback: "~"},             // nf-fa-pencil
	IconGitUntracked:      {Glyph: "\uF059", Fallback: "?"},             // nf-fa-question
	IconGitAhead:          {Glyph: "↑", Fallback: "↑"},                  // plain unicode
	IconGitBehind:         {Glyph: "↓", Fallback: "↓"},                  // plain unicode
	IconGitCleanDot:       {Glyph: "●", Fallback: "●"},                  // plain unicode
	IconPR:                {Glyph: "\uF407", Fallback: "PR"},            // nf-oct-git_pull_request
	IconPRDraft:           {Glyph: "\uF5D9", Fallback: "[draft]"},       // nf-oct-git_pull_request_draft (approx)
	IconPRPending:         {Glyph: "\U000F0DA5", Fallback: "[review]"},  // nf-md-clock_alert_outline (approx)
	IconPRApproved:        {Glyph: "\uF42E", Fallback: "[ok]"},          // nf-oct-check
	IconPRChangesRequest:  {Glyph: "\U000F0159", Fallback: "[!]"},       // nf-md-close_circle
	IconContextWindow:     {Glyph: "\U000F035B", Fallback: "ctx:"},      // nf-md-memory
	IconContextAlert:      {Glyph: "\U000F0026", Fallback: "!!"},        // nf-md-alert
	IconCost:              {Glyph: "\uF155", Fallback: "$"},             // nf-fa-dollar
	IconDuration:          {Glyph: "\uF017", Fallback: "⏱"},             // nf-fa-clock_o
	IconRateLimitFiveHour: {Glyph: "\U000F051F", Fallback: "5h:"},       // nf-md-timer_sand
	IconRateLimitWeek:     {Glyph: "\U000F0F94", Fallback: "7d:"},       // nf-md-calendar_week
	IconVim:               {Glyph: "\uF11C", Fallback: "[mode]"},        // nf-fa-keyboard_o
	IconAgent:             {Glyph: "\U000F0904", Fallback: "agent:"},    // nf-md-account_hard_hat
	IconOutputStyle:       {Glyph: "\U000F0765", Fallback: "style:"},    // nf-md-palette_outline
	IconThinking:          {Glyph: "\uF0EB", Fallback: "thinking"},      // nf-fa-lightbulb_o
	IconCache:             {Glyph: "\U000F0224", Fallback: "cache"},     // nf-md-database_outline (approx)
	IconSessionName:       {Glyph: "\uF02B", Fallback: "session:"},      // nf-fa-tag
	IconRepo:              {Glyph: "\uF401", Fallback: "repo:"},         // nf-oct-repo
	IconWorktree:          {Glyph: "\U000F0928", Fallback: "worktree:"}, // nf-md-source_branch (approx)
}

// Glyph returns the icon for key: the Nerd Font glyph when nerdFont is true,
// otherwise the plain-text fallback. Returns "" for an unknown key.
func Glyph(key string, nerdFont bool) string {
	icon, ok := Icons[key]
	if !ok {
		return ""
	}
	if nerdFont {
		return icon.Glyph
	}
	return icon.Fallback
}
