package type_filter

type Phone string

func (p Phone) Sanitize() Phone {
	if len(p) <= 4 {
		return Phone("****")
	}
	return Phone("***" + string(p[len(p)-4:]))
}

type Secret string

func (s Secret) Sanitize() Secret {
	return Secret("***")
}

type ContactInfo struct {
	Phone Phone
}

type User struct {
	Name    string
	Contact ContactInfo
}

// Unrelated has a maskable field but should NOT be generated
// when using -type=User.
type Unrelated struct {
	Token Secret
}
