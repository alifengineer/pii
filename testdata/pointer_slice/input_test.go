package pointer_slice

import "testing"

func TestUser_Sanitize_SliceMask(t *testing.T) {
	u := User{
		Name:   "Alice",
		Phones: []Phone{Phone("+1234567890"), Phone("+0987654321")},
	}

	sanitized := u.Sanitize()

	if sanitized.Phones[0] != Phone("***7890") {
		t.Errorf("Phones[0] = %q, want %q", sanitized.Phones[0], "***7890")
	}
	if sanitized.Phones[1] != Phone("***4321") {
		t.Errorf("Phones[1] = %q, want %q", sanitized.Phones[1], "***4321")
	}
	if u.Phones[0] != Phone("+1234567890") {
		t.Error("original slice should not be mutated")
	}
}

func TestUser_Sanitize_NilPointer(t *testing.T) {
	u := User{
		Name:    "Alice",
		Address: nil,
	}
	sanitized := u.Sanitize()
	if sanitized.Address != nil {
		t.Error("nil pointer should stay nil")
	}
}

func TestUser_Sanitize_NonNilPointer(t *testing.T) {
	addr := &Address{
		Street: "123 Main St",
		Zip:    Phone("12345"),
	}
	u := User{
		Name:    "Alice",
		Address: addr,
	}

	sanitized := u.Sanitize()

	if sanitized.Address == nil {
		t.Fatal("non-nil pointer should not become nil")
	}
	if sanitized.Address.Zip != Phone("***2345") {
		t.Errorf("Address.Zip = %q, want %q", sanitized.Address.Zip, "***2345")
	}
	if sanitized.Address.Street != "123 Main St" {
		t.Errorf("Address.Street = %q, want %q", sanitized.Address.Street, "123 Main St")
	}
	if addr.Zip != Phone("12345") {
		t.Error("original should not be mutated")
	}
}
