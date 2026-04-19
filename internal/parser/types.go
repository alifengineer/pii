package parser

type StructInfo struct {
	Name   string
	Fields []FieldInfo
}

type FieldInfo struct {
	Name        string
	TypeName    string
	HasSanitize bool
	IsStruct    bool
	IsPtr       bool
	IsSlice     bool
	ElemType    string
}
