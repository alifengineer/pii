package basic

type Phone string

func (p Phone) Sanitize() Phone {
	if len(p) <= 4 {
		return Phone("****")
	}
	return Phone("***" + string(p[len(p)-4:]))
}

type User struct {
	Name  string
	Phone Phone
}
