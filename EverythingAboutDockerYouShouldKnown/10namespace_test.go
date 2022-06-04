package namespace

import "testing"

func TestPidMode_Valid(t *testing.T) {
	tests := []struct {
		name string
		n    PidMode
		want bool
	}{
		// TODO: Add test cases.
		{"base case", "", true},
		{"host case", "host", true},
		{"container case", "container", false},
		{"container case2", "container:1", true},
		{"default case", "con", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Valid(); got != tt.want {
				t.Errorf("PidMode.Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
