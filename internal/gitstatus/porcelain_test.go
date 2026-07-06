package gitstatus

import (
	"fmt"
	"strings"
	"testing"
)

func TestParsePorcelainV2(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want Status
	}{
		{
			name: "clean tree, no upstream",
			in: "# branch.oid abc1234\n" +
				"# branch.head main\n",
			want: Status{Branch: "main", OID: "abc1234"},
		},
		{
			name: "detached head",
			in: "# branch.oid abc1234\n" +
				"# branch.head (detached)\n",
			want: Status{OID: "abc1234", Detached: true},
		},
		{
			name: "ahead and behind upstream",
			in: "# branch.oid abc1234\n" +
				"# branch.head main\n" +
				"# branch.upstream origin/main\n" +
				"# branch.ab +2 -3\n",
			want: Status{Branch: "main", OID: "abc1234", Upstream: "origin/main", Ahead: 2, Behind: 3},
		},
		{
			name: "staged and modified files",
			in: "# branch.head main\n" +
				"1 M. N... 100644 100644 100644 hash1 hash2 staged_only.go\n" +
				"1 .M N... 100644 100644 100644 hash1 hash2 modified_only.go\n" +
				"1 MM N... 100644 100644 100644 hash1 hash2 both.go\n",
			want: Status{Branch: "main", Staged: 2, Modified: 2},
		},
		{
			name: "renamed file counts via type 2 entries",
			in: "# branch.head main\n" +
				"2 R. N... 100644 100644 100644 hash1 hash2 R100 new.go\told.go\n",
			want: Status{Branch: "main", Staged: 1},
		},
		{
			name: "untracked files",
			in: "# branch.head main\n" +
				"? new_file.go\n" +
				"? another.go\n",
			want: Status{Branch: "main", Untracked: 2},
		},
		{
			name: "unmerged conflict",
			in: "# branch.head main\n" +
				"u UU N... 100644 100644 100644 100644 hash1 hash2 hash3 conflict.go\n",
			want: Status{Branch: "main", Conflicts: 1},
		},
		{
			name: "empty output",
			in:   "",
			want: Status{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParsePorcelainV2([]byte(tt.in))
			if err != nil {
				t.Fatalf("ParsePorcelainV2() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("ParsePorcelainV2() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestStatusClean(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		st   Status
		want bool
	}{
		{name: "all zero is clean", st: Status{}, want: true},
		{name: "staged makes it dirty", st: Status{Staged: 1}, want: false},
		{name: "modified makes it dirty", st: Status{Modified: 1}, want: false},
		{name: "untracked makes it dirty", st: Status{Untracked: 1}, want: false},
		{name: "conflicts make it dirty", st: Status{Conflicts: 1}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.st.Clean(); got != tt.want {
				t.Errorf("Clean() = %v, want %v", got, tt.want)
			}
		})
	}
}

// largePorcelainFixture builds a realistic worst-case porcelain-v2 output: a
// branch header plus a large mixed-status working tree, closer to a big
// monorepo mid-refactor than the tiny few-line fixtures above.
func largePorcelainFixture(fileCount int) []byte {
	var b strings.Builder
	b.WriteString("# branch.oid abc1234567890\n")
	b.WriteString("# branch.head main\n")
	b.WriteString("# branch.upstream origin/main\n")
	b.WriteString("# branch.ab +3 -5\n")
	for i := range fileCount {
		switch i % 4 {
		case 0:
			fmt.Fprintf(&b, "1 M. N... 100644 100644 100644 hash1 hash2 pkg/staged%d.go\n", i)
		case 1:
			fmt.Fprintf(&b, "1 .M N... 100644 100644 100644 hash1 hash2 pkg/modified%d.go\n", i)
		case 2:
			fmt.Fprintf(&b, "? pkg/untracked%d.go\n", i)
		case 3:
			fmt.Fprintf(&b, "2 R. N... 100644 100644 100644 hash1 hash2 R100 pkg/renamed%d.go\tpkg/old%d.go\n", i, i)
		}
	}
	return []byte(b.String())
}

// BenchmarkParsePorcelainV2 measures the pure parser against a realistic
// large working tree (500 changed files), the one piece of gitstatus that
// doesn't depend on an external git subprocess and so is meaningful to
// benchmark in isolation.
func BenchmarkParsePorcelainV2(b *testing.B) {
	fixture := largePorcelainFixture(500)
	for b.Loop() {
		if _, err := ParsePorcelainV2(fixture); err != nil {
			b.Fatalf("ParsePorcelainV2() error = %v", err)
		}
	}
}
