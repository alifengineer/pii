package sanitize

import (
	"strings"
	"unicode/utf8"
)

type Name string

func (n Name) Sanitize() Name {
	s := string(n)
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
	return Name(strings.Join(masked, " "))
}
