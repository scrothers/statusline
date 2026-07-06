// Command statusline is a Claude Code statusLine command: it reads session
// JSON from stdin and prints a themed, Nerd-Font-powerline status line to
// stdout. See https://code.claude.com/docs/en/statusline.
package main

import "os"

func main() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
