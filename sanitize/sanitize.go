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
