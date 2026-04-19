package analyzer

type FieldAction int

const (
	ActionCopy FieldAction = iota
	ActionSanitize
	ActionSanitizePtr
	ActionSanitizeSlice
	ActionSanitizePtrSlice
)

type SanitizeTarget struct {
	StructName string
	Fields     []TargetField
}

type TargetField struct {
	Name   string
	Action FieldAction
}
