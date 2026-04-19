package mask

import (
	"strings"
	"unicode/utf8"
)

// MaskName keeps the first rune, stars the rest. "Jonathan" -> "J*******".
func MaskName[T ~string](v T) T {
	n := utf8.RuneCountInString(string(v))
	switch n {
	case 0:
		return T("")
	case 1:
		return T("*")
	}

	r, _ := utf8.DecodeRuneInString(string(v))
	return T(string(r) + strings.Repeat("*", n-1))
}

// MaskPhone keeps the last 4 digits of the numeric portion.
// "+998 (90) 123-45-67" -> "********4567".
func MaskPhone[T ~string](v T) T {
	if v == "" {
		return T("")
	}
	if len(v) <= 4 {
		return T(strings.Repeat("*", len(v)))
	}

	return T(strings.Repeat("*", len(v)-4) + string(v)[len(v)-4:])
}

// MaskEmail keeps the first char of the local part and the full domain.
// "john.doe@acme.com" -> "j*******@acme.com".
func MaskEmail[T ~string](v T) T {
	at := strings.LastIndex(string(v), "@")
	if at <= 0 {
		return T(strings.Repeat("*", utf8.RuneCountInString(string(v))))
	}

	local, domain := string(v)[:at], string(v)[at:]
	n := utf8.RuneCountInString(local)
	if n <= 1 {
		return T(local + domain)
	}

	r, _ := utf8.DecodeRuneInString(local)
	return T(string(r) + strings.Repeat("*", n-1) + domain)
}

// MaskPAN keeps the first 6 (BIN) and last 4 digits — PCI-DSS safe.
// "4111 1111 1111 1111" -> "411111******1111".
func MaskPAN[T ~string](v T) T {
	return T(string(v)[:6] + strings.Repeat("*", len(v)-10) + string(v)[len(v)-4:])
}
