package pii

type TagValue string

const (
	Name  TagValue = "name"
	Phone TagValue = "phone"
	Email TagValue = "email"
	PAN   TagValue = "pan"
)

const piiTag = "pii"

type Masker func(value string) string

var maskers = map[TagValue]Masker{
	Name:  MaskName[string],
	Phone: MaskPhone[string],
	Email: MaskEmail[string],
	PAN:   MaskPAN[string],
}

func Register(tag TagValue, m Masker) {
	maskers[tag] = m
}
