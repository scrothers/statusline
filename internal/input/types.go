package input

// Payload is the full JSON object Claude Code pipes to the statusLine
// command's stdin on every refresh. Every field the docs mark as "absent"
// or "nullable" under some condition is a pointer here, so a Go zero value
// never gets confused with "not present yet".
type Payload struct {
	CWD            string         `json:"cwd"`
	SessionID      string         `json:"session_id"`
	SessionName    string         `json:"session_name,omitempty"`
	PromptID       string         `json:"prompt_id,omitempty"`
	TranscriptPath string         `json:"transcript_path"`
	Version        string         `json:"version"`
	Model          *Model         `json:"model"`
	Workspace      *Workspace     `json:"workspace"`
	OutputStyle    *OutputStyle   `json:"output_style"`
	Cost           *Cost          `json:"cost"`
	ContextWindow  *ContextWindow `json:"context_window"`
	Exceeds200k    bool           `json:"exceeds_200k_tokens"`
	Effort         *Effort        `json:"effort"`
	Thinking       *Thinking      `json:"thinking"`
	RateLimits     *RateLimits    `json:"rate_limits"`
	Vim            *Vim           `json:"vim"`
	Agent          *Agent         `json:"agent"`
	PR             *PR            `json:"pr"`
	Worktree       *Worktree      `json:"worktree"`
}

// Model identifies the current model backing the session.
type Model struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

// Workspace describes the directories and repository context for the session.
type Workspace struct {
	CurrentDir  string   `json:"current_dir"`
	ProjectDir  string   `json:"project_dir"`
	AddedDirs   []string `json:"added_dirs"`
	GitWorktree string   `json:"git_worktree,omitempty"`
	Repo        *Repo    `json:"repo"`
}

// Repo identifies the git remote origin, parsed by Claude Code from the
// workspace's origin remote. Absent outside a git repository or when no
// origin remote is configured.
type Repo struct {
	Host  string `json:"host"`
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

// OutputStyle names the current output style.
type OutputStyle struct {
	Name string `json:"name"`
}

// Cost carries session cost, duration, and line-change accounting.
type Cost struct {
	TotalCostUSD       float64 `json:"total_cost_usd"`
	TotalDurationMS    int64   `json:"total_duration_ms"`
	TotalAPIDurationMS int64   `json:"total_api_duration_ms"`
	TotalLinesAdded    int     `json:"total_lines_added"`
	TotalLinesRemoved  int     `json:"total_lines_removed"`
}

// ContextWindow describes context-window usage from the most recent API
// response. UsedPercentage, RemainingPercentage, and CurrentUsage may be nil
// early in a session or immediately after /compact.
type ContextWindow struct {
	TotalInputTokens    int      `json:"total_input_tokens"`
	TotalOutputTokens   int      `json:"total_output_tokens"`
	ContextWindowSize   int      `json:"context_window_size"`
	UsedPercentage      *float64 `json:"used_percentage"`
	RemainingPercentage *float64 `json:"remaining_percentage"`
	CurrentUsage        *Usage   `json:"current_usage"`
}

// Usage breaks ContextWindow's token totals out by category.
type Usage struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
}

// Effort reports the current reasoning effort level. Absent when the
// current model does not support the effort parameter.
type Effort struct {
	Level string `json:"level"`
}

// Thinking reports whether extended thinking is enabled for the session.
type Thinking struct {
	Enabled bool `json:"enabled"`
}

// RateLimits carries Claude subscription rate-limit usage. Present only for
// Claude.ai subscribers after the first API response in the session; each
// window may be independently absent.
type RateLimits struct {
	FiveHour *RateLimitWindow `json:"five_hour"`
	SevenDay *RateLimitWindow `json:"seven_day"`
}

// RateLimitWindow reports usage and reset time for one rate-limit window.
type RateLimitWindow struct {
	UsedPercentage float64 `json:"used_percentage"`
	ResetsAt       int64   `json:"resets_at"`
}

// Vim reports the current vim mode. Present only when vim mode is enabled.
type Vim struct {
	Mode string `json:"mode"`
}

// Agent identifies the running agent. Present only when running with the
// --agent flag or agent settings configured.
type Agent struct {
	Name string `json:"name"`
}

// PR describes the open pull request for the current branch, if any.
// ReviewState may be independently absent even when PR is present.
type PR struct {
	Number      int    `json:"number"`
	URL         string `json:"url"`
	ReviewState string `json:"review_state,omitempty"`
}

// Worktree describes an active --worktree session. Branch and
// OriginalBranch may be absent for hook-based worktrees.
type Worktree struct {
	Name           string `json:"name"`
	Path           string `json:"path"`
	Branch         string `json:"branch,omitempty"`
	OriginalCWD    string `json:"original_cwd"`
	OriginalBranch string `json:"original_branch,omitempty"`
}

// Deref returns *p, or def if p is nil. It exists to keep segment code free
// of repetitive nil-check boilerplate for simple scalar fallbacks; deeply
// nested optional paths (Workspace.Repo, ContextWindow.CurrentUsage) still
// get an explicit nil check at the point of use.
func Deref[T any](p *T, def T) T {
	if p == nil {
		return def
	}
	return *p
}
