// Package render composes segment output into the final statusline text: it
// joins each line's segments with powerline cap/connector glyphs (or a
// plain divider between pill-less badges), drops lowest-priority segments
// under width pressure, and never lets one segment's panic take down the
// whole line.
package render
