package type_filter

import "testing"

func TestUser_Sanitize(t *testing.T) {
	u := User{
		Name: "Alice",
		Contact: ContactInfo{
			Phone: Phone("+1234567890"),
		},
	}

	s := u.Sanitize()

	if s.Name != "Alice" {
		t.Errorf("Name = %q, want %q", s.Name, "Alice")
	}
	if s.Contact.Phone != Phone("***7890") {
		t.Errorf("Contact.Phone = %q, want %q", s.Contact.Phone, "***7890")
	}
}

func TestContactInfo_Sanitize(t *testing.T) {
	c := ContactInfo{
		Phone: Phone("+1234567890"),
	}

	s := c.Sanitize()

	if s.Phone != Phone("***7890") {
		t.Errorf("Phone = %q, want %q", s.Phone, "***7890")
	}
}
