# Contributing to statusline

Thanks for considering a contribution. This is a small, single-binary Go
CLI, so the process is scoped accordingly.

## Development setup

```sh
git clone https://github.com/scrothers/statusline.git
cd statusline
make build
make test
```

See the [README's Development section](README.md#development) for the full
command list: `make test-integration`, `make test-e2e`, `make bench`,
`make lint`, `make security`.

## Before opening a pull request

- `make fmt lint test test-integration test-e2e security` all clean.
- New segments, themes, or icons: verify any new Nerd Font glyph against
  the authoritative
  [glyphnames.json](https://raw.githubusercontent.com/ryanoasis/nerd-fonts/master/glyphnames.json)
  rather than guessing a codepoint from memory — codepoints have drifted
  across Nerd Font releases even when names stayed stable.
- Keep commits atomic, with imperative-mood subjects (`segment: add X`, not
  `Added X`).
- Add or update tests for any behavior change. `internal/segment`,
  `internal/render`, and `internal/style` currently sit at 100% coverage —
  new code there should keep it that way.
- No AI/co-author trailers in commit messages.

## Reporting bugs or requesting features

Open an issue using the provided templates. For security vulnerabilities,
see [SECURITY.md](SECURITY.md) instead of filing a public issue.

## Code of conduct

This project follows the [Contributor Covenant](CODE_OF_CONDUCT.md).
