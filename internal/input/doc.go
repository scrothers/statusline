// Package input defines the JSON schema Claude Code sends to a statusLine
// command on stdin and parses it defensively, tolerating absent, null, and
// unknown fields so the binary never fails on partial or future data.
package input
