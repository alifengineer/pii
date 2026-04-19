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
