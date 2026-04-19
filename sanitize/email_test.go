package sanitize

import "testing"

func TestEmail_Sanitize(t *testing.T) {
	tests := []struct {
		in   Email
		want Email
	}{
		{"alice@example.com", "a***@e***.com"},
		{"bob@mail.co.uk", "b***@m***.uk"},
		{"x@y.io", "x***@y***.io"},
		{"noatsign", "****"},
		{"", ""},
		{"a@b", "a***@b***"},
		{"@empty.com", "***@e***.com"},
		{"local@", "l***@***"},
	}
	for _, tt := range tests {
		got := tt.in.Sanitize()
		if got != tt.want {
			t.Errorf("Email(%q).Sanitize() = %q, want %q", tt.in, got, tt.want)
		}
	}
}
