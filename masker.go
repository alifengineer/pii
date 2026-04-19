package pii

import (
	"strings"
	"unicode/utf8"
)

// MaskName keeps the first rune, stars the rest. "Jonathan" -> "J*******".
func MaskName[T ~string](v T) T {
	s := string(v)
	words := strings.Fields(s)
	if len(words) == 0 {
		return ""
	}

	masked := make([]string, len(words))
	for i, w := range words {
		first, size := utf8.DecodeRuneInString(w)
		if size == 0 {
			masked[i] = w
			continue
		}
		remaining := utf8.RuneCountInString(w) - 1
		masked[i] = string(first) + strings.Repeat("*", remaining)
	}

	return T(strings.Join(masked, " "))
}

// MaskPhone keeps the last 4 digits of the numeric portion.
// "+998 (90) 123-45-67" -> "********4567".
func MaskPhone[T ~string](v T) T {
	if v == "" {
		return ""
	}

	if len(v) <= 4 {
		return "****"
	}

	return T("***" + string(v[len(v)-4:]))
}

// MaskEmail keeps the first char of the local part and the full domain.
// "john.doe@acme.com" -> "j*******@acme.com".
func MaskEmail[T ~string](v T) T {
	s := string(v)
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
		return T(maskedLocal + "@" + maskFirstLetter(domain))
	}

	domainPrefix := domain[:dot]
	tld := domain[dot+1:]
	return T(maskedLocal + "@" + maskFirstLetter(domainPrefix) + "." + tld)
}

func maskFirstLetter(s string) string {
	if s == "" {
		return "***"
	}
	return string(s[0]) + "***"
}

// MaskPAN keeps the first 6 (BIN) and last 4 digits — PCI-DSS safe.
// "4111 1111 1111 1111" -> "411111******1111".
func MaskPAN[T ~string](v T) T {
	s := string(v)
	if s == "" {
		return ""
	}
	if len(s) <= 10 {
		return "****"
	}

	middle := len(s) - 10
	return T(s[:6] + strings.Repeat("*", middle) + s[len(s)-4:])
}
