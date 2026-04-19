package analyzer

import (
	"testing"

	"github.com/alifengineer/pii/internal/parser"
)

func TestAnalyze_BasicMask(t *testing.T) {
	structs := []parser.StructInfo{
		{
			Name: "User",
			Fields: []parser.FieldInfo{
				{Name: "Name", TypeName: "string"},
				{Name: "Phone", TypeName: "Phone", HasSanitize: true, ElemType: "Phone"},
			},
		},
	}

	targets, _, err := Analyze(structs, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 1 {
		t.Fatalf("got %d targets, want 1", len(targets))
	}
	if targets[0].StructName != "User" {
		t.Errorf("target = %q, want %q", targets[0].StructName, "User")
	}

	assertFieldAction(t, targets[0], "Name", ActionCopy)
	assertFieldAction(t, targets[0], "Phone", ActionSanitize)
}

func TestAnalyze_NestedStructs(t *testing.T) {
	structs := []parser.StructInfo{
		{
			Name: "User",
			Fields: []parser.FieldInfo{
				{Name: "Name", TypeName: "string"},
				{Name: "ContactInfo", TypeName: "ContactInfo", IsStruct: true, ElemType: "ContactInfo"},
			},
		},
		{
			Name: "ContactInfo",
			Fields: []parser.FieldInfo{
				{Name: "Phone", TypeName: "Phone", HasSanitize: true, ElemType: "Phone"},
			},
		},
	}

	targets, _, err := Analyze(structs, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 2 {
		t.Fatalf("got %d targets, want 2", len(targets))
	}
	if targets[0].StructName != "ContactInfo" {
		t.Errorf("first target = %q, want %q", targets[0].StructName, "ContactInfo")
	}
	if targets[1].StructName != "User" {
		t.Errorf("second target = %q, want %q", targets[1].StructName, "User")
	}
}

func TestAnalyze_DeeplyNested(t *testing.T) {
	structs := []parser.StructInfo{
		{Name: "A", Fields: []parser.FieldInfo{
			{Name: "B", IsStruct: true, ElemType: "B"},
		}},
		{Name: "B", Fields: []parser.FieldInfo{
			{Name: "C", IsStruct: true, ElemType: "C"},
		}},
		{Name: "C", Fields: []parser.FieldInfo{
			{Name: "D", IsStruct: true, ElemType: "D"},
		}},
		{Name: "D", Fields: []parser.FieldInfo{
			{Name: "Val", HasSanitize: true, ElemType: "Val"},
		}},
	}

	targets, _, err := Analyze(structs, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 4 {
		t.Fatalf("got %d targets, want 4", len(targets))
	}

	want := []string{"D", "C", "B", "A"}
	for i, w := range want {
		if targets[i].StructName != w {
			t.Errorf("targets[%d] = %q, want %q", i, targets[i].StructName, w)
		}
	}
}

func TestAnalyze_NoMaskableFields(t *testing.T) {
	structs := []parser.StructInfo{
		{
			Name: "Config",
			Fields: []parser.FieldInfo{
				{Name: "Host", TypeName: "string"},
				{Name: "Port", TypeName: "int"},
			},
		},
	}

	targets, _, err := Analyze(structs, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 0 {
		t.Errorf("got %d targets, want 0", len(targets))
	}
}

func TestAnalyze_PointerField(t *testing.T) {
	structs := []parser.StructInfo{
		{
			Name: "Wrapper",
			Fields: []parser.FieldInfo{
				{Name: "Addr", IsStruct: true, IsPtr: true, ElemType: "Addr"},
			},
		},
		{
			Name: "Addr",
			Fields: []parser.FieldInfo{
				{Name: "Zip", HasSanitize: true, ElemType: "Zip"},
			},
		},
	}

	targets, _, err := Analyze(structs, nil)
	if err != nil {
		t.Fatal(err)
	}

	wrapper := findTarget(t, targets, "Wrapper")
	assertFieldAction(t, wrapper, "Addr", ActionSanitizePtr)
}

func TestAnalyze_SliceField(t *testing.T) {
	structs := []parser.StructInfo{
		{
			Name: "Wrapper",
			Fields: []parser.FieldInfo{
				{Name: "Phones", IsSlice: true, HasSanitize: true, ElemType: "Phone"},
			},
		},
	}

	targets, _, err := Analyze(structs, nil)
	if err != nil {
		t.Fatal(err)
	}

	wrapper := findTarget(t, targets, "Wrapper")
	assertFieldAction(t, wrapper, "Phones", ActionSanitizeSlice)
}

func TestAnalyze_MixedActions(t *testing.T) {
	structs := []parser.StructInfo{
		{
			Name: "Mixed",
			Fields: []parser.FieldInfo{
				{Name: "Name", TypeName: "string"},
				{Name: "Phone", HasSanitize: true, ElemType: "Phone"},
				{Name: "Info", IsStruct: true, ElemType: "Info"},
			},
		},
		{
			Name: "Info",
			Fields: []parser.FieldInfo{
				{Name: "Email", HasSanitize: true, ElemType: "Email"},
			},
		},
	}

	targets, _, err := Analyze(structs, nil)
	if err != nil {
		t.Fatal(err)
	}

	mixed := findTarget(t, targets, "Mixed")
	assertFieldAction(t, mixed, "Name", ActionCopy)
	assertFieldAction(t, mixed, "Phone", ActionSanitize)
	assertFieldAction(t, mixed, "Info", ActionSanitize)
}

func TestAnalyze_TypeFilter_SingleType(t *testing.T) {
	structs := []parser.StructInfo{
		{
			Name: "User",
			Fields: []parser.FieldInfo{
				{Name: "Phone", HasSanitize: true, ElemType: "Phone"},
			},
		},
		{
			Name: "Order",
			Fields: []parser.FieldInfo{
				{Name: "CardNum", HasSanitize: true, ElemType: "CardNum"},
			},
		},
	}

	targets, _, err := Analyze(structs, map[string]bool{"User": true})
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 1 {
		t.Fatalf("got %d targets, want 1", len(targets))
	}
	if targets[0].StructName != "User" {
		t.Errorf("target = %q, want %q", targets[0].StructName, "User")
	}
}

func TestAnalyze_TypeFilter_MultipleTypes(t *testing.T) {
	structs := []parser.StructInfo{
		{
			Name: "User",
			Fields: []parser.FieldInfo{
				{Name: "Phone", HasSanitize: true, ElemType: "Phone"},
			},
		},
		{
			Name: "Order",
			Fields: []parser.FieldInfo{
				{Name: "CardNum", HasSanitize: true, ElemType: "CardNum"},
			},
		},
		{
			Name: "Internal",
			Fields: []parser.FieldInfo{
				{Name: "Secret", HasSanitize: true, ElemType: "Secret"},
			},
		},
	}

	targets, _, err := Analyze(structs, map[string]bool{"User": true, "Order": true})
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 2 {
		t.Fatalf("got %d targets, want 2", len(targets))
	}
	findTarget(t, targets, "User")
	findTarget(t, targets, "Order")
}

func TestAnalyze_TypeFilter_TransitiveDep(t *testing.T) {
	structs := []parser.StructInfo{
		{
			Name: "User",
			Fields: []parser.FieldInfo{
				{Name: "Name", TypeName: "string"},
				{Name: "Contact", IsStruct: true, ElemType: "ContactInfo"},
			},
		},
		{
			Name: "ContactInfo",
			Fields: []parser.FieldInfo{
				{Name: "Phone", HasSanitize: true, ElemType: "Phone"},
			},
		},
		{
			Name: "Unrelated",
			Fields: []parser.FieldInfo{
				{Name: "Secret", HasSanitize: true, ElemType: "Secret"},
			},
		},
	}

	targets, autoIncluded, err := Analyze(structs, map[string]bool{"User": true})
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 2 {
		t.Fatalf("got %d targets, want 2", len(targets))
	}
	if targets[0].StructName != "ContactInfo" {
		t.Errorf("first target = %q, want %q", targets[0].StructName, "ContactInfo")
	}
	if targets[1].StructName != "User" {
		t.Errorf("second target = %q, want %q", targets[1].StructName, "User")
	}
	if len(autoIncluded) != 1 || autoIncluded[0] != "ContactInfo" {
		t.Errorf("autoIncluded = %v, want [ContactInfo]", autoIncluded)
	}
}

func TestAnalyze_TypeFilter_DeepTransitiveDep(t *testing.T) {
	structs := []parser.StructInfo{
		{Name: "User", Fields: []parser.FieldInfo{
			{Name: "Contact", IsStruct: true, ElemType: "ContactInfo"},
		}},
		{Name: "ContactInfo", Fields: []parser.FieldInfo{
			{Name: "Addr", IsStruct: true, ElemType: "Address"},
		}},
		{Name: "Address", Fields: []parser.FieldInfo{
			{Name: "Zip", HasSanitize: true, ElemType: "Zip"},
		}},
		{Name: "Unrelated", Fields: []parser.FieldInfo{
			{Name: "Secret", HasSanitize: true, ElemType: "Secret"},
		}},
	}

	targets, _, err := Analyze(structs, map[string]bool{"User": true})
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 3 {
		t.Fatalf("got %d targets, want 3", len(targets))
	}
	want := []string{"Address", "ContactInfo", "User"}
	for i, w := range want {
		if targets[i].StructName != w {
			t.Errorf("targets[%d] = %q, want %q", i, targets[i].StructName, w)
		}
	}
}

func TestAnalyze_TypeFilter_DepOnly(t *testing.T) {
	structs := []parser.StructInfo{
		{
			Name: "User",
			Fields: []parser.FieldInfo{
				{Name: "Contact", IsStruct: true, ElemType: "ContactInfo"},
			},
		},
		{
			Name: "ContactInfo",
			Fields: []parser.FieldInfo{
				{Name: "Phone", HasSanitize: true, ElemType: "Phone"},
			},
		},
	}

	targets, _, err := Analyze(structs, map[string]bool{"ContactInfo": true})
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 1 {
		t.Fatalf("got %d targets, want 1", len(targets))
	}
	if targets[0].StructName != "ContactInfo" {
		t.Errorf("target = %q, want %q", targets[0].StructName, "ContactInfo")
	}
}

func TestAnalyze_TypeFilter_NoMaskExplicit(t *testing.T) {
	structs := []parser.StructInfo{
		{
			Name: "Config",
			Fields: []parser.FieldInfo{
				{Name: "Host", TypeName: "string"},
				{Name: "Port", TypeName: "int"},
			},
		},
	}

	targets, _, err := Analyze(structs, map[string]bool{"Config": true})
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 1 {
		t.Fatalf("got %d targets, want 1", len(targets))
	}
	if targets[0].StructName != "Config" {
		t.Errorf("target = %q, want %q", targets[0].StructName, "Config")
	}
	assertFieldAction(t, targets[0], "Host", ActionCopy)
	assertFieldAction(t, targets[0], "Port", ActionCopy)
}

func TestAnalyze_TypeFilter_Nonexistent(t *testing.T) {
	structs := []parser.StructInfo{
		{
			Name: "User",
			Fields: []parser.FieldInfo{
				{Name: "Phone", HasSanitize: true, ElemType: "Phone"},
			},
		},
	}

	_, _, err := Analyze(structs, map[string]bool{"Nonexistent": true})
	if err == nil {
		t.Fatal("expected error for nonexistent type")
	}
}

func TestAnalyze_TypeFilter_MixedValidInvalid(t *testing.T) {
	structs := []parser.StructInfo{
		{
			Name: "User",
			Fields: []parser.FieldInfo{
				{Name: "Phone", HasSanitize: true, ElemType: "Phone"},
			},
		},
	}

	_, _, err := Analyze(structs, map[string]bool{"User": true, "Nonexistent": true})
	if err == nil {
		t.Fatal("expected error when any type name is invalid")
	}
}

func TestAnalyze_AutoIncluded_BothExplicit(t *testing.T) {
	structs := []parser.StructInfo{
		{
			Name: "User",
			Fields: []parser.FieldInfo{
				{Name: "Contact", IsStruct: true, ElemType: "ContactInfo"},
			},
		},
		{
			Name: "ContactInfo",
			Fields: []parser.FieldInfo{
				{Name: "Phone", HasSanitize: true, ElemType: "Phone"},
			},
		},
	}

	_, autoIncluded, err := Analyze(structs, map[string]bool{"User": true, "ContactInfo": true})
	if err != nil {
		t.Fatal(err)
	}
	if len(autoIncluded) != 0 {
		t.Errorf("autoIncluded = %v, want empty", autoIncluded)
	}
}

func TestAnalyze_AutoIncluded_NilFilter(t *testing.T) {
	structs := []parser.StructInfo{
		{
			Name: "User",
			Fields: []parser.FieldInfo{
				{Name: "Phone", HasSanitize: true, ElemType: "Phone"},
			},
		},
	}

	_, autoIncluded, err := Analyze(structs, nil)
	if err != nil {
		t.Fatal(err)
	}
	if autoIncluded != nil {
		t.Errorf("autoIncluded = %v, want nil", autoIncluded)
	}
}

func TestClassifyField_AllActions(t *testing.T) {
	sanitizable := map[string]bool{"Inner": true}

	tests := []struct {
		name  string
		field parser.FieldInfo
		want  FieldAction
	}{
		{"copy", parser.FieldInfo{Name: "X"}, ActionCopy},
		{"sanitize_leaf", parser.FieldInfo{Name: "X", HasSanitize: true, ElemType: "T"}, ActionSanitize},
		{"sanitize_slice", parser.FieldInfo{Name: "X", HasSanitize: true, IsSlice: true, ElemType: "T"}, ActionSanitizeSlice},
		{"sanitize", parser.FieldInfo{Name: "X", IsStruct: true, ElemType: "Inner"}, ActionSanitize},
		{"sanitize_ptr", parser.FieldInfo{Name: "X", IsStruct: true, IsPtr: true, ElemType: "Inner"}, ActionSanitizePtr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := classifyField(tt.field, sanitizable)
			if got != tt.want {
				t.Errorf("classifyField = %d, want %d", got, tt.want)
			}
		})
	}
}

func findTarget(t *testing.T, targets []SanitizeTarget, name string) SanitizeTarget {
	t.Helper()
	for _, tgt := range targets {
		if tgt.StructName == name {
			return tgt
		}
	}
	t.Fatalf("target %q not found", name)
	return SanitizeTarget{}
}

func assertFieldAction(t *testing.T, target SanitizeTarget, fieldName string, want FieldAction) {
	t.Helper()
	for _, f := range target.Fields {
		if f.Name == fieldName {
			if f.Action != want {
				t.Errorf("%s.%s action = %d, want %d", target.StructName, fieldName, f.Action, want)
			}
			return
		}
	}
	t.Errorf("field %q not found in target %q", fieldName, target.StructName)
}
