// Package config loads the statusline's TOML configuration: a built-in
// default merged with an optional user file, discovered by a fixed
// precedence and never failing hard — a missing or malformed config always
// degrades to the built-in default plus a warning.
package config
