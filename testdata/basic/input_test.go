package basic

import "testing"

func TestUser_Sanitize(t *testing.T) {
	u := User{
		Name:  "Alice",
		Phone: Phone("+1234567890"),
	}

	sanitized := u.Sanitize()

	if sanitized.Phone == u.Phone {
		t.Error("Phone should be masked")
	}
	if sanitized.Phone != Phone("***7890") {
		t.Errorf("Phone = %q, want %q", sanitized.Phone, "***7890")
	}
	if sanitized.Name != "Alice" {
		t.Errorf("Name = %q, want %q", sanitized.Name, "Alice")
	}
	if u.Phone != Phone("+1234567890") {
		t.Error("original should not be mutated")
	}
}
