package mask_types

import (
	"testing"

	"github.com/alifengineer/pii/sanitize"
)

func TestCustomer_Sanitize(t *testing.T) {
	c := Customer{
		ID:    42,
		Name:  sanitize.Name("Alice Smith"),
		Phone: sanitize.Phone("+1234567890"),
		Email: sanitize.Email("alice@example.com"),
		Card:  sanitize.PAN("4111111111111111"),
	}

	s := c.Sanitize()

	if s.ID != 42 {
		t.Errorf("ID = %d, want 42", s.ID)
	}
	if s.Name != sanitize.Name("A**** S****") {
		t.Errorf("Name = %q, want %q", s.Name, "A**** S****")
	}
	if s.Phone != sanitize.Phone("***7890") {
		t.Errorf("Phone = %q, want %q", s.Phone, "***7890")
	}
	if s.Email != sanitize.Email("a***@e***.com") {
		t.Errorf("Email = %q, want %q", s.Email, "a***@e***.com")
	}
	if s.Card != sanitize.PAN("411111******1111") {
		t.Errorf("Card = %q, want %q", s.Card, "411111******1111")
	}

	// Original not mutated
	if c.Phone != sanitize.Phone("+1234567890") {
		t.Error("original should not be mutated")
	}
}
