package sanitize

import "testing"

func TestPAN_Sanitize(t *testing.T) {
	tests := []struct {
		in   PAN
		want PAN
	}{
		{"4111111111111111", "411111******1111"},
		{"1234567890", "****"},
		{"12345678901", "123456*8901"},
		{"", ""},
		{"1234567890123456789", "123456*********6789"},
		{"12345678901234", "123456****1234"},
	}
	for _, tt := range tests {
		got := tt.in.Sanitize()
		if got != tt.want {
			t.Errorf("PAN(%q).Sanitize() = %q, want %q", tt.in, got, tt.want)
		}
	}
}
