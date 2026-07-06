# statusline

A single-binary, themeable [Claude Code statusLine](https://code.claude.com/docs/en/statusline)
command written in Go. Flat, no-background segments joined by a Nerd Font
divider, truecolor, five built-in themes, and an optional TOML config file
for anyone who wants to tweak it further.

```
[38;2;254;128;25m󰚩[0m[38;2;213;196;161m Opus[0m[38;2;146;131;116m    [0m[38;2;250;189;47m[0m[38;2;146;131;116m    [0m[38;2;241;196;15m󰊚 high[0m[38;2;146;131;116m    [0m[38;2;250;147;48m󰍛[0m[38;2;146;131;116m ⟨[0m[38;2;184;187;38m█[0m[38;2;190;187;38m█[0m[38;2;197;187;39m█[0m[38;2;204;187;40m█[0m[38;2;211;187;41m█[0m[38;2;218;188;42m█[0m[38;2;225;188;43m█[0m[38;2;232;188;44m█[0m[38;2;239;188;45m█[0m[38;2;246;188;46m█[0m[38;2;250;182;47m█[0m[38;2;250;170;47m█[0m[38;2;250;158;48m█[0m[38;2;250;146;48m▋[0m[38;2;80;73;69m░░░░░░[0m[38;2;146;131;116m⟩ [0m[38;2;250;147;48m68%[0m[38;2;146;131;116m 136.0k/64.0k[0m[38;2;146;131;116m    [0m[38;2;239;188;45m 5h[0m[38;2;146;131;116m ⟨[0m[38;2;184;187;38m█[0m[38;2;210;187;41m█[0m[38;2;236;188;45m▌[0m[38;2;80;73;69m░░░[0m[38;2;146;131;116m⟩[0m[38;2;239;188;45m 42%[0m[38;2;146;131;116m  2h[0m[38;2;146;131;116m    [0m[38;2;250;140;49m󰨳 7d[0m[38;2;146;131;116m ⟨[0m[38;2;184;187;38m█[0m[38;2;210;187;41m█[0m[38;2;236;188;45m█[0m[38;2;250;165;48m█[0m[38;2;250;119;50m▎[0m[38;2;80;73;69m░[0m[38;2;146;131;116m⟩[0m[38;2;250;140;49m 71%[0m[38;2;146;131;116m  3d 2h[0m[38;2;146;131;116m    [0m[38;2;211;187;41m[0m[38;2;211;187;41m 79%[0m[38;2;146;131;116m (108.0k)[0m
[38;2;254;128;25m[0m[38;2;213;196;161m big-refactor[0m[38;2;146;131;116m    [0m[38;2;254;128;25m󰉋[0m[38;2;213;196;161m /home/user/code/statusline[0m[38;2;146;131;116m    [0m[38;2;250;189;47m  [0m[38;2;184;187;38m[0m[38;2;213;196;161m 342[0m[38;2;251;73;52m [0m[38;2;213;196;161m 58[0m[38;2;146;131;116m    [0m[38;2;250;189;47m  [0m[38;2;184;187;38m󰄠[0m[38;2;213;196;161m 20.0k[0m[38;2;251;73;52m 󰄝[0m[38;2;213;196;161m 4.5k[0m[38;2;142;192;124m 󰆺[0m[38;2;213;196;161m 8.0k[0m[38;2;142;192;124m [0m[38;2;213;196;161m 108.0k[0m[38;2;146;131;116m    [0m[38;2;184;187;38m[0m[38;2;235;219;178m2.17[0m[38;2;146;131;116m    [0m[38;2;142;192;124m[0m[38;2;235;219;178m 1:23:45[0m
[38;2;254;128;25m[0m[38;2;146;131;116m github.com[0m[38;2;213;196;161m/scrothers[0m[38;2;213;196;161m/statusline[0m[38;2;146;131;116m    [0m[38;2;184;187;38m[0m[38;2;213;196;161m #128[0m[38;2;184;187;38m approved[0m[38;2;146;131;116m    [0m[38;2;250;189;47m[0m[38;2;213;196;161m main[0m[38;2;184;187;38m 󰔡[0m[38;2;213;196;161m 2[0m[38;2;250;189;47m  [0m[38;2;213;196;161m 1[0m[38;2;146;131;116m [0m[38;2;213;196;161m 3[0m[38;2;184;187;38m ↑[0m[38;2;213;196;161m 1[0m[38;2;146;131;116m    [0m[38;2;254;128;25m󰤨[0m[38;2;213;196;161m my-feature[0m
```

(a real render — `statusline demo --theme gruvbox --scenario full`)

## Install

From source:

```sh
git clone https://github.com/scrothers/statusline.git
cd statusline
make build      # produces ./statusline
```

Or, once a release exists:

```sh
go install github.com/scrothers/statusline/cmd/statusline@latest
```

## Preview it

`statusline demo` renders built-in sample payloads directly — no need to
craft JSON fixtures or wire up Claude Code first:

```sh
./statusline demo                                  # "full" scenario, all 5 themes
./statusline demo --theme dracula                   # one theme only
./statusline demo --theme nord --scenario minimal    # early-session look
./statusline demo --scenario narrow --columns 40     # test width-based truncation
```

Scenarios: `minimal` (early session, no git repo yet), `full` (every segment
at once — dirty repo, open PR, context/cost/rate-limit/cache data, a named
session, a worktree), `narrow` (the `full` payload rendered at 30 columns, to
see which segments drop first).

## Configure Claude Code

Add to `~/.claude/settings.json`:

```json
{
  "statusLine": {
    "type": "command",
    "command": "~/code/statusline/statusline",
    "padding": 2,
    "refreshInterval": 2,
    "hideVimModeIndicator": true
  }
}
```

- `refreshInterval: 2` keeps the clock-style duration, cost, and context bar
  ticking during long idle stretches (a slow tool call, a background
  subagent), not just after each assistant message.
- `hideVimModeIndicator: true` suppresses Claude Code's built-in
  `-- INSERT --` text, since the vim badge already shows the mode when
  enabled in config (see [Segments](#segments)).

## Segments

Three lines, each answering one question. Every segment is plain colored
text — no background is ever painted, so lines read as flat, breathing
text joined by a thin chevron divider ( `` ), not powerline pills.

| Line | Segments | Notes |
|---|---|---|
| 1 — Claude | model, thinking, effort, context window, rate limits (5h/7d), cache | The model's own state: what it's running, how hard, and how much room/budget is left. |
| 2 — session | session name, directory, lines added/removed, token counts, cost, duration | Omitted fields (no custom name, no diff yet) just don't appear. |
| 3 — git | repo (host/owner/name), open PR (number + review state), branch + status, worktree | The whole line disappears outside a git repository. |

Segments not in the default layout but still available via custom config:
`vim` (vim mode), `agent` (subagent name), `output_style`. These render with
no background either — they just aren't wired into any line by default.

Under width pressure, segments drop in priority order (lowest first) until
a line fits; model and directory never drop, only self-truncate.

`effort` is the one segment colored by intensity rather than theme: its icon
escalates from an empty gauge (`low`) through a full gauge (`xhigh`) to fire
(`max`) and an alert fire (`ultra`), and its color runs a fixed green → red →
purple scale independent of the active theme, so "getting hotter" reads the
same everywhere.

`context_window` and the two `ratelimit_*` gauges share the same bar
treatment: `context_window`'s width scales with the detected terminal width
(clamped between 8 and 24 cells); the rate-limit bars stay a fixed, narrower
6 cells and are explicitly labeled `5h`/`7d` so the two aren't
distinguishable by icon alone. On all three, each bar cell's color is fixed
by its position along the bar — green on the left sliding through warning
to danger on the right — so filling the bar reveals more of a stable
on-screen gradient from the left rather than shifting every already-filled
cell's color together each time the percentage changes. (Each gauge's icon
and percentage text use a separate smooth gradient based on its own overall
percentage.) The context bar also shows a `used/remaining` token count next
to the percentage whenever the context window size is known.

Each `ratelimit_*` gauge also shows a muted reset countdown (a restart icon
plus a compact "resets in" duration like `2h` or `3d 2h`) whenever the
payload reports a reset time — a coarser, more actionable signal than the
percentage alone, since it tells you when the budget actually comes back
rather than just how much is used right now.

`cache` reuses that same smooth gradient for its icon and hit-rate
percentage, but inverted: a high cache-hit rate is good, so it runs green at
100% down to red at 0% (the opposite direction from the gauges above, where
high usage is the thing to watch). The raw cache-read count in parentheses
stays a plain muted color.

`lines_changed` leads with a pencil icon (given two trailing spaces of its
own, since its glyph reads tight against a following icon), then
diff-added/diff-removed icons that alone carry the +/- meaning. `token_counts`
leads with a ticket icon (same two-space treatment) and breaks down the most
recent API response's usage into four counts, each behind its own icon: an
inbound tray for input tokens, an outbound tray for output tokens, a
database-plus for cache-creation tokens, and the same cache icon used
elsewhere for cache-read tokens. All four come from that one response, so
they share a single time scope rather than mixing session totals with a
per-response cache breakdown. In both segments, only the icons carry
semantic color (green/red for add-remove and input-output, info for the
cache pair); the counts themselves are always the theme's secondary text
color, with no ASCII sign.

`cost` and `duration` follow the same icon/text split: the dollar icon is
success-green (money) and the clock icon is info, while the amount and the
clock-style duration both render in the theme's primary text color — plain
prose, not a semantic accent, since a dollar figure or elapsed time isn't
inherently good/bad/informational the way a gauge fill is.

Line 3 follows the same rules as everywhere else. `git`'s status badges
(staged/modified/untracked/conflicts/ahead/behind) each split into a
category-colored icon and a `TextSecondary` count, matching
`lines_changed`/`token_counts`; conflicts use a proper alert icon rather
than a bare `!`. The branch name, `repo`'s final `/name` piece, `pr`'s
number, and `worktree`'s name all use `IdentityText` — the same "headline
label following an identity-colored icon" role as `model`/`directory`/
`session_name` — while `pr`'s icon and review-state word still share the
review-state color, since that state is the actual signal.

`repo`'s icon is host-branded: GitHub, GitLab, and Forgejo each get their
own logo when the remote host contains that name (matching public domains
like `github.com` and enterprise/self-hosted ones like
`github.company.com` alike), falling back to a generic git icon for
anything else — including Forgejo instances with no "forgejo" in their
hostname (e.g. Codeberg), since the host string is the only signal
available.

## Themes

Seven built-in themes, selected with `theme = "<name>"` in config or
`--theme <name>` on the command line. `gruvbox` is the default. Themes are
foreground-only palettes (identity accent + success/warning/danger/info/muted
roles) — there's no background token, since the statusline never paints one.

| Name | Style |
|---|---|
| `gruvbox` | Warm, retro, high-contrast (default) |
| `catppuccin-mocha` | Soft pastel-on-dark |
| `tokyo-night` | Cool blues/purples on near-black |
| `nord` | Arctic blues |
| `dracula` | High-contrast purple/pink |
| `claude-dark` | Anthropic's coral/terracotta accent on a warm near-black ink background |
| `claude-light` | The same coral/terracotta accent on a warm cream background — the one built-in theme meant for a light terminal |

`claude-dark` and `claude-light` are built from Anthropic's published brand
palette (the `#D97757` coral, `#141413` ink, `#FAF9F5` cream, plus a sage
green and dusty blue used for `success`/`info`); `warning` and `danger`
aren't part of that palette, so they're original colors chosen to sit
comfortably in the same warm, muted family rather than clashing with a
generic bright red/yellow.

## Configuration

Optional TOML config file, discovered in this order (first match wins):

1. `--config <path>` flag
2. `$XDG_CONFIG_HOME/statusline/config.toml` (or `~/.config/statusline/config.toml`)
3. `~/.claude/statusline-config.toml`
4. Built-in defaults (Gruvbox theme, the 3-line layout above)

A missing or malformed config file is never fatal — it falls back to the
built-in default and prints a warning to stderr (never stdout, which is
reserved for the rendered statusline).

Example overriding the theme and disabling one segment:

```toml
theme = "nord"

[segments.pr]
enabled = false

[theme_overrides]
success = "#00ff00"
```

Take over `lines[].segments` entirely to reorder, drop, or add segments —
see the [Segments](#segments) table for every available ID:

```toml
[[lines]]
enabled = true
segments = ["model", "context_window", "cache"]
```

## Requirements

- A [Nerd Font](https://www.nerdfonts.com/) in your terminal for the icons
  and the divider glyph. Without one, set `nerd_font = false` in config to
  fall back to plain Unicode/ASCII for every icon.
- `git` on `PATH` for the repo/branch/PR segments; everything else works
  without it.

## Development

```sh
make build              # go build -> ./statusline
make test               # unit tests
make test-integration   # + real git subprocess tests
make test-e2e           # + builds and drives the real binary
make bench               # benchmarks for the render/parse/config/theme hot paths
make lint                # go vet + golangci-lint
make security            # govulncheck + gosec
```

The whole non-git-subprocess render path — JSON decode, config/theme load,
and the full three-line render — benchmarks at roughly 1.3ms combined on
modest hardware, comfortably inside the sub-100ms-per-invocation budget the
design targets; `internal/gitstatus`'s porcelain parser alone handles a
500-file working tree in well under a millisecond.

## License

Apache License 2.0 — see [LICENSE](LICENSE).
