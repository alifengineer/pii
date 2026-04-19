package sanitize

type Phone string

func (p Phone) Sanitize() Phone {
	if len(p) == 0 {
		return ""
	}
	if len(p) <= 4 {
		return "****"
	}
	return Phone("***" + string(p[len(p)-4:]))
}
