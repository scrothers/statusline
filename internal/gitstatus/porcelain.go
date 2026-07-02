package gitstatus

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Status is the git branch and working-tree state for one directory.
type Status struct {
	Branch    string `json:"branch"`
	OID       string `json:"oid"`
	Detached  bool   `json:"detached"`
	Upstream  string `json:"upstream,omitempty"`
	Ahead     int    `json:"ahead"`
	Behind    int    `json:"behind"`
	Staged    int    `json:"staged"`
	Modified  int    `json:"modified"`
	Untracked int    `json:"untracked"`
	Conflicts int    `json:"conflicts"`
	NotARepo  bool   `json:"not_a_repo"`
}

// Clean reports whether the working tree has no staged, modified,
// untracked, or conflicted files.
func (s Status) Clean() bool {
	return s.Staged == 0 && s.Modified == 0 && s.Untracked == 0 && s.Conflicts == 0
}

// ParsePorcelainV2 parses the output of
// `git status --porcelain=v2 --branch --untracked-files=normal`. It's a
// pure function so it's unit-testable against fixture bytes without ever
// invoking git.
func ParsePorcelainV2(out []byte) (Status, error) {
	var st Status

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "#":
			parseBranchHeader(&st, fields)
		case "1", "2":
			if len(fields) < 2 || len(fields[1]) != 2 {
				continue
			}
			xy := fields[1]
			if xy[0] != '.' {
				st.Staged++
			}
			if xy[1] != '.' {
				st.Modified++
			}
		case "u":
			st.Conflicts++
		case "?":
			st.Untracked++
		}
	}
	if err := scanner.Err(); err != nil {
		return Status{}, fmt.Errorf("gitstatus: parse porcelain v2: %w", err)
	}
	return st, nil
}

func parseBranchHeader(st *Status, fields []string) {
	if len(fields) < 3 {
		return
	}
	switch fields[1] {
	case "branch.oid":
		st.OID = fields[2]
	case "branch.head":
		if fields[2] == "(detached)" {
			st.Detached = true
		} else {
			st.Branch = fields[2]
		}
	case "branch.upstream":
		st.Upstream = fields[2]
	case "branch.ab":
		if len(fields) < 4 {
			return
		}
		st.Ahead = parseSignedCount(fields[2])
		st.Behind = parseSignedCount(fields[3])
	}
}

// parseSignedCount parses a branch.ab count field like "+3" or "-0",
// returning 0 for anything that doesn't parse rather than erroring — an
// unrecognized ahead/behind field shouldn't fail the entire status parse.
func parseSignedCount(s string) int {
	n, err := strconv.Atoi(strings.TrimPrefix(s, "+"))
	if err != nil {
		return 0
	}
	if n < 0 {
		return -n
	}
	return n
}
