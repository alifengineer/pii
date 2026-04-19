package mask_types

import "github.com/alifengineer/pii/sanitize"

type Customer struct {
	ID    int
	Name  sanitize.Name
	Phone sanitize.Phone
	Email sanitize.Email
	Card  sanitize.PAN
}
