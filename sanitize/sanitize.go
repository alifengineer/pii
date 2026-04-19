// Package sanitize provides built-in sanitizable types for use with the pii code generator.
//
// Each type is a named string type with a Sanitize() method that returns a sanitized
// copy. When used as struct fields, pii's generator automatically detects the
// Sanitize() method and generates Sanitize() methods that call it.
//
// Usage:
//
//	import "github.com/alifengineer/pii/sanitize"
//
//	type User struct {
//	    Name  sanitize.Name
//	    Phone sanitize.Phone
//	    Email sanitize.Email
//	}
//
// Then run: //go:generate pii
package sanitize

import "github.com/alifengineer/pii"

type Email string

func (e Email) Sanitize() Email {
	return pii.MaskEmail(e)
}

type Name string

func (n Name) Sanitize() Name {
	return pii.MaskName(n)
}

type PAN string

func (p PAN) Sanitize() PAN {
	return pii.MaskPAN(p)
}

type Phone string

func (p Phone) Sanitize() Phone {
	return pii.MaskPhone(p)
}
