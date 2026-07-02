package config

import (
	_ "embed"
	"fmt"

	"github.com/BurntSushi/toml"
)

//go:embed default.toml
var defaultTOML string

// Default returns the built-in default configuration. It never fails: a
// broken embedded default.toml would be a build-time bug, not a runtime
// condition, so a decode error here panics rather than propagating a
// spurious runtime error path that could never actually be exercised by a
// user-supplied file.
func Default() Config {
	var cfg Config
	if _, err := toml.Decode(defaultTOML, &cfg); err != nil {
		panic(fmt.Sprintf("config: embedded default.toml is invalid: %v", err))
	}
	return cfg
}
