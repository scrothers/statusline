package render

import "testing"

func TestColumns(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  int
	}{
		{name: "empty falls back to default", value: "", want: defaultColumns},
		{name: "valid value is used", value: "120", want: 120},
		{name: "non-numeric falls back to default", value: "not-a-number", want: defaultColumns},
		{name: "zero falls back to default", value: "0", want: defaultColumns},
		{name: "negative falls back to default", value: "-10", want: defaultColumns},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("COLUMNS", tt.value)
			if got := Columns(); got != tt.want {
				t.Errorf("Columns() = %d, want %d", got, tt.want)
			}
		})
	}
}
