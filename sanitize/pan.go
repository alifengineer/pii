package sanitize

import "strings"

type PAN string

func (p PAN) Sanitize() PAN {
	s := string(p)
	if s == "" {
		return ""
	}
	if len(s) <= 10 {
		return "****"
	}
	middle := len(s) - 10
	return PAN(s[:6] + strings.Repeat("*", middle) + s[len(s)-4:])
}
