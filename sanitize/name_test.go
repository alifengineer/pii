package sanitize

import "testing"

func TestName_Sanitize(t *testing.T) {
	tests := []struct {
		in   Name
		want Name
	}{
		{"Alice", "A****"},
		{"Bob Smith", "B** S****"},
		{"Mary-Jane", "M********"},
		{"A", "A"},
		{"", ""},
		{"  ", ""},
		{"Алиса", "А****"},
		{"A B C", "A B C"},
	}
	for _, tt := range tests {
		got := tt.in.Sanitize()
		if got != tt.want {
			t.Errorf("Name(%q).Sanitize() = %q, want %q", tt.in, got, tt.want)
		}
	}
}
