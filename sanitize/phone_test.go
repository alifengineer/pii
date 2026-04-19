package sanitize

import "testing"

func TestPhone_Sanitize(t *testing.T) {
	tests := []struct {
		in   Phone
		want Phone
	}{
		{"+1234567890", "***7890"},
		{"12345", "***2345"},
		{"1234", "****"},
		{"123", "****"},
		{"1", "****"},
		{"", ""},
	}
	for _, tt := range tests {
		got := tt.in.Sanitize()
		if got != tt.want {
			t.Errorf("Phone(%q).Sanitize() = %q, want %q", tt.in, got, tt.want)
		}
	}
}
