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
	IconThinkingOn        = "thinking_on"
	IconThinkingOff       = "thinking_off"
	IconCache             = "cache"
	IconSessionName       = "session_name"
	IconRepoGitHub        = "repo_github"
	IconRepoGitLab        = "repo_gitlab"
	IconRepoForgejo       = "repo_forgejo"
	IconRepoGit           = "repo_git"
	IconWorktree          = "worktree"
	IconEffortLow         = "effort_low"
	IconEffortMedium      = "effort_medium"
	IconEffortHigh        = "effort_high"
	IconEffortXHigh       = "effort_xhigh"
	IconEffortMax         = "effort_max"
	IconEffortUltra       = "effort_ultra"
	IconLinesAdded        = "lines_added"
	IconLinesRemoved      = "lines_removed"
	IconEditPencil        = "edit_pencil"
	IconTokensTicket      = "tokens_ticket"
	IconTokensInput       = "tokens_input"
	IconTokensOutput      = "tokens_output"
	IconTokensCacheCreate = "tokens_cache_create"
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
	IconPR:                {Glyph: "\uF407", Fallback: "PR"},            // nf-oct-git_pull_request
	IconPRDraft:           {Glyph: "\uF5D9", Fallback: "[draft]"},       // nf-oct-git_pull_request_draft (approx)
	IconPRPending:         {Glyph: "\U000F0DA5", Fallback: "[review]"},  // nf-md-clock_alert_outline (approx)
	IconPRApproved:        {Glyph: "\uF42E", Fallback: "[ok]"},          // nf-oct-check
	IconPRChangesRequest:  {Glyph: "\U000F0159", Fallback: "[!]"},       // nf-md-close_circle
	IconContextWindow:     {Glyph: "\U000F035B", Fallback: "ctx:"},      // nf-md-memory
	IconContextAlert:      {Glyph: "\U000F0026", Fallback: "!!"},        // nf-md-alert
	IconCost:              {Glyph: "\uF155", Fallback: "$"},             // nf-fa-dollar
	IconDuration:          {Glyph: "\uF017", Fallback: "⏱"},             // nf-fa-clock_o
	IconRateLimitFiveHour: {Glyph: "\uF252", Fallback: "5h:"},           // nf-fa-hourglass_2
	IconRateLimitWeek:     {Glyph: "\U000F0A33", Fallback: "7d:"},       // nf-md-calendar_week
	IconVim:               {Glyph: "\uF11C", Fallback: "[mode]"},        // nf-fa-keyboard_o
	IconAgent:             {Glyph: "\U000F0904", Fallback: "agent:"},    // nf-md-account_hard_hat
	IconOutputStyle:       {Glyph: "\U000F0765", Fallback: "style:"},    // nf-md-palette_outline
	IconThinkingOn:        {Glyph: "\uEA61", Fallback: "thinking"},      // nf-cod-lightbulb
	IconThinkingOff:       {Glyph: "\uF0EB", Fallback: "idle"},          // nf-fa-lightbulb
	IconCache:             {Glyph: "\uF49B", Fallback: "cache"},         // nf-oct-cache
	IconSessionName:       {Glyph: "\uF02B", Fallback: "session:"},      // nf-fa-tag
	IconRepoGitHub:        {Glyph: "\uF09B", Fallback: "GH:"},           // nf-fa-github
	IconRepoGitLab:        {Glyph: "\uE7EB", Fallback: "GL:"},           // nf-dev-gitlab
	IconRepoForgejo:       {Glyph: "\uF335", Fallback: "FJ:"},           // nf-linux-forgejo
	IconRepoGit:           {Glyph: "\uE702", Fallback: "repo:"},         // nf-dev-git
	IconWorktree:          {Glyph: "\U000F0928", Fallback: "worktree:"}, // nf-md-source_branch (approx)
	IconEffortLow:         {Glyph: "\U000F0873", Fallback: "low"},       // nf-md-gauge_empty
	IconEffortMedium:      {Glyph: "\U000F0875", Fallback: "medium"},    // nf-md-gauge_low
	IconEffortHigh:        {Glyph: "\U000F029A", Fallback: "high"},      // nf-md-gauge
	IconEffortXHigh:       {Glyph: "\U000F0874", Fallback: "xhigh"},     // nf-md-gauge_full
	IconEffortMax:         {Glyph: "\U000F0238", Fallback: "max"},       // nf-md-fire
	IconEffortUltra:       {Glyph: "\U000F15D7", Fallback: "ultra"},     // nf-md-fire_alert
	IconLinesAdded:        {Glyph: "\uF457", Fallback: "+"},             // nf-oct-diff_added
	IconLinesRemoved:      {Glyph: "\uF458", Fallback: "-"},             // nf-oct-diff_removed
	IconEditPencil:        {Glyph: "\uF040", Fallback: "*"},             // nf-fa-pencil (same glyph as git_modified)
	IconTokensTicket:      {Glyph: "\uF145", Fallback: "tok:"},          // nf-fa-ticket
	IconTokensInput:       {Glyph: "\U000F0120", Fallback: "in"},        // nf-md-tray_arrow_down
	IconTokensOutput:      {Glyph: "\U000F011D", Fallback: "out"},       // nf-md-tray_arrow_up
	IconTokensCacheCreate: {Glyph: "\U000F01BA", Fallback: "c+"},        // nf-md-database_plus
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
