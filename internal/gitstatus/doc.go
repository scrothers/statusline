// Package gitstatus collects git branch and working-tree status for the git
// segment, caching results per session so repeated statusline refreshes
// don't re-run git on every trigger.
package gitstatus
