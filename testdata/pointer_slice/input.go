package pointer_slice

type Phone string

func (p Phone) Sanitize() Phone {
	if len(p) <= 4 {
		return Phone("****")
	}
	return Phone("***" + string(p[len(p)-4:]))
}

type Address struct {
	Street string
	Zip    Phone
}

type User struct {
	Name    string
	Phones  []Phone
	Address *Address
}
