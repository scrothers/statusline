// Package claudetheme detects Claude Code's own active color theme —
// "dark" or "light" — by reading Claude Code's own settings.json files
// directly, so statusline's palette always matches what Claude Code itself
// is rendering with. It never fails hard: any missing, unreadable, or
// unrecognized setting degrades to "dark", the same default Claude Code
// itself documents.
package claudetheme
