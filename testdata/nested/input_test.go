package nested

import "testing"

func TestContactInfo_Sanitize(t *testing.T) {
	c := ContactInfo{
		Phone: Phone("+1234567890"),
		Email: Email("alice@example.com"),
	}

	sanitized := c.Sanitize()

	if sanitized.Phone != Phone("***7890") {
		t.Errorf("Phone = %q, want %q", sanitized.Phone, "***7890")
	}
	if sanitized.Email != Email("***@***.***") {
		t.Errorf("Email = %q, want %q", sanitized.Email, "***@***.***")
	}
}

func TestUser_Sanitize(t *testing.T) {
	u := User{
		Name: "Alice",
		ContactInfo: ContactInfo{
			Phone: Phone("+1234567890"),
			Email: Email("alice@example.com"),
		},
		Address: Address{
			Street: "123 Main St",
			Zip:    "12345",
		},
	}

	sanitized := u.Sanitize()

	if sanitized.ContactInfo.Phone != Phone("***7890") {
		t.Errorf("ContactInfo.Phone = %q, want %q", sanitized.ContactInfo.Phone, "***7890")
	}
	if sanitized.ContactInfo.Email != Email("***@***.***") {
		t.Errorf("ContactInfo.Email = %q, want %q", sanitized.ContactInfo.Email, "***@***.***")
	}
	if sanitized.Name != "Alice" {
		t.Errorf("Name = %q, want %q", sanitized.Name, "Alice")
	}
	if sanitized.Address.Street != "123 Main St" {
		t.Errorf("Address.Street = %q, want %q", sanitized.Address.Street, "123 Main St")
	}
}
