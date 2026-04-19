package sanitize

import "strings"

type Email string

func (e Email) Sanitize() Email {
	s := string(e)
	if s == "" {
		return ""
	}

	before, after, ok := strings.Cut(s, "@")
	if !ok {
		return "****"
	}

	local := before
	domain := after

	maskedLocal := maskFirstLetter(local)

	dot := strings.LastIndexByte(domain, '.')
	if dot < 0 {
		return Email(maskedLocal + "@" + maskFirstLetter(domain))
	}

	domainPrefix := domain[:dot]
	tld := domain[dot+1:]

	return Email(maskedLocal + "@" + maskFirstLetter(domainPrefix) + "." + tld)
}

func maskFirstLetter(s string) string {
	if s == "" {
		return "***"
	}
	return string(s[0]) + "***"
}
