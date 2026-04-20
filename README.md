[![CI](https://github.com/alifengineer/pii/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/alifengineer/pii/actions/workflows/ci.yml)
[![cov](https://alifengineer.github.io/pii/badges/coverage.svg)](https://github.com/alifengineer/pii/actions)

# pii — PII Masking & Data Sanitization Library for Go

Mask sensitive fields in Go structs at runtime or compile time. Supports names, emails, phones, and PANs with custom masker registration. 
This covers the long-tail queries: "mask struct fields golang", "pii redaction go library", "sanitize personal data go".

Pii generates `Sanitize()` methods for Go structs. Any field with a `Sanitize()` method gets sanitized automatically; everything else is copied as-is. The result is a clean, non-mutating sanitizer you can call before logging or serializing sensitive data.

## Install

```sh
go install github.com/alifengineer/pii@latest
```

## Quick Start

Define a struct using sanitize types:

```go
package user

import "github.com/alifengineer/pii/sanitize"

//go:generate pii

type User struct {
    ID    int
    Name  sanitize.Name
    Phone sanitize.Phone
    Email sanitize.Email
}
```

Run `go generate`:

```sh
go generate ./...
```

Pii creates `input_sanitize.go` with:

```go
func (v User) Sanitize() User {
    v.Name = v.Name.Sanitize()
    v.Phone = v.Phone.Sanitize()
    v.Email = v.Email.Sanitize()
    return v
}
```

Use it:

```go
sanitized := user.Sanitize()
log.Println(sanitized) // {42 A**** S**** ***7890 a***@e***.com}
```

## Built-in Sanitize Types

| Type | Example Input | Sanitized Output |
|------|---------------|------------------|
| `sanitize.Phone` | `+1234567890` | `***7890` |
| `sanitize.Email` | `alice@example.com` | `a***@e***.com` |
| `sanitize.Name` | `Alice Smith` | `A**** S****` |
| `sanitize.PAN` | `4111111111111111` | `411111******1111` |

## Custom Sanitize Types

Any type with a `Sanitize()` method that returns itself works:

```go
type SSN string

func (s SSN) Sanitize() SSN {
    if len(s) < 4 {
        return "****"
    }
    return SSN("***-**-" + string(s[len(s)-4:]))
}
```

Pii detects `Sanitize()` on struct fields and calls it in the generated `Sanitize()` method.

## Selective Generation

Use `-type` to generate for specific structs only:

```go
//go:generate pii -type=User,Order
```

When `-type` is set:

- Only named structs get `Sanitize()` methods
- Transitive dependencies are auto-included (if `User` has a `ContactInfo` field that needs sanitizing, `ContactInfo.Sanitize()` is generated too)
- Structs not named and not reachable as dependencies are skipped
- Unknown type names produce a clear error

When `-type` is omitted, all qualifying structs in the package are generated (default behavior).

## Nested Structs

Pii handles struct fields that themselves contain sanitizable fields:

```go
type ContactInfo struct {
    Phone sanitize.Phone
    Email sanitize.Email
}

type User struct {
    Name    string
    Contact ContactInfo
}
```

Both `ContactInfo.Sanitize()` and `User.Sanitize()` are generated. `User.Sanitize()` calls `Contact.Sanitize()` internally.

Pointer fields (`*ContactInfo`) and slices are also handled.

## Flags

| Flag | Description |
|------|-------------|
| `-type=Name1,Name2` | Generate only for named structs and their dependencies |
| `-destination=path` | Output file path (default: `<source>_sanitize.go`) |
| `-check` | Compare output to existing file, exit non-zero if different |
| `-v` | Print per-field reasoning and auto-included dependencies to stderr |
| `-debug` | Dump parsed model to stdout without generating |

## License

[MIT](LICENSE)
