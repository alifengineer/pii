package nested

type Phone string

func (p Phone) Sanitize() Phone {
	if len(p) <= 4 {
		return Phone("****")
	}
	return Phone("***" + string(p[len(p)-4:]))
}

type Email string

func (e Email) Sanitize() Email {
	return Email("***@***.***")
}

type Address struct {
	Street string
	Zip    string
}

type ContactInfo struct {
	Phone Phone
	Email Email
}

//go:generate sanitizer -type=User
type User struct {
	Name        string
	ContactInfo ContactInfo
	Address     Address
}
