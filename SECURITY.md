# Security Policy

## Supported versions

statusline ships as a single binary with no maintained older release
branches — only the latest release receives security fixes.

| Version  | Supported |
| -------- | --------- |
| Latest   | ✅        |
| Older    | ❌        |

## Reporting a vulnerability

Please **do not** open a public issue for a security vulnerability.

Use GitHub's private vulnerability reporting for this repository:
[Report a vulnerability](https://github.com/scrothers/statusline/security/advisories/new)
(also reachable from the Security tab). This opens a private advisory with
the maintainer, keeps the report confidential until a fix ships, and lets
you request a CVE through GitHub if warranted.

Expect an initial response within a few days. Once a fix is available,
we'll agree on a disclosure timeline with you before any public advisory is
published.

## Scope

statusline reads JSON from stdin (from Claude Code), an optional local TOML
config file, and shells out to `git status` in the current working
directory. Relevant vulnerability classes include:

- Command injection via the git subprocess invocation
- Path traversal in config or cache file handling
- Denial of service via a malformed or oversized stdin payload
- Unsafe behavior triggered by a malicious config file

It makes no network requests, so network-based attack vectors are out of
scope.
