// Package fallback provides the last-resort degraded statusline used when
// normal rendering can't proceed at all (stdin parse failure, an
// unrecovered panic that reached main, or an empty render result), so the
// binary always prints something rather than leaving the status bar blank.
package fallback
