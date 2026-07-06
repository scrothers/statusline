# statusline

A single-binary, themeable [Claude Code statusLine](https://code.claude.com/docs/en/statusline)
command written in Go. Powerline-style segments, Nerd Font icons, truecolor,
five built-in themes, and an optional TOML config file for anyone who wants
to tweak it further.

```
 󰚩 Opus  󰉋 ~/code/statusline 
  master 12   #128 
 󰍛 ⟨██████▊░░░⟩ 68%  $0.42 1:23:45  󰔟██▌░░░ 42%  󰾔████▎░ 71%  INSERT · 󰤄 reviewer · high
```

(a real render — captured from this project's own working tree)

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
at once — dirty repo, open PR, context/cost/rate-limit data, all bonus
badges), `narrow` (the `full` payload rendered at 30 columns, to see which
segments drop first).

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
  `-- INSERT --` text, since the vim badge already shows the mode.

## Segments

| Line | Segments | Notes |
|---|---|---|
| 1 — identity | model, directory | Always shown; directory breadcrumb-truncates under width pressure. |
| 2 — repo | git branch/status, open PR | Omitted entirely outside a git repository. |
| 3 — vitals | context window, cost + duration, rate limits (5h/7d), bonus badges (vim mode, agent name, effort level, output style) | Bonus badges are the first thing dropped in a narrow terminal. |

## Themes

Five built-in themes, selected with `theme = "<name>"` in config or
`--theme <name>` on the command line. `gruvbox` is the default.

| Name | Style |
|---|---|
| `gruvbox` | Warm, retro, high-contrast (default) |
| `catppuccin-mocha` | Soft pastel-on-dark |
| `tokyo-night` | Cool blues/purples on near-black |
| `nord` | Arctic blues |
| `dracula` | High-contrast purple/pink |

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

## Requirements

- A [Nerd Font](https://www.nerdfonts.com/) in your terminal for the icons
  and powerline separators. Without one, set `nerd_font = false` in config
  to fall back to plain Unicode/ASCII for every icon.
- `git` on `PATH` for the repository segment; everything else works without it.

## Development

```sh
make build              # go build -> ./statusline
make test               # unit tests
make test-integration   # + real git subprocess tests
make test-e2e           # + builds and drives the real binary
make lint                # go vet + golangci-lint
```

## License

Apache License 2.0 — see [LICENSE](LICENSE).
