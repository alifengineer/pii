package parser

import (
	"path/filepath"
	"runtime"
	"testing"
)

func testdataDir(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", name)
}

func TestParse_BasicStruct(t *testing.T) {
	pkg, structs, err := Parse(testdataDir("basic"))
	if err != nil {
		t.Fatal(err)
	}
	if pkg != "basic" {
		t.Errorf("package = %q, want %q", pkg, "basic")
	}

	user := findStruct(t, structs, "User")
	if len(user.Fields) != 2 {
		t.Fatalf("User has %d fields, want 2", len(user.Fields))
	}

	nameField := findField(t, user, "Name")
	if nameField.HasSanitize {
		t.Error("Name should not have Sanitize()")
	}

	phoneField := findField(t, user, "Phone")
	if !phoneField.HasSanitize {
		t.Error("Phone should have Sanitize()")
	}
	if phoneField.IsPtr || phoneField.IsSlice || phoneField.IsStruct {
		t.Error("Phone should not be ptr/slice/struct")
	}
}

func TestParse_NestedStructs(t *testing.T) {
	_, structs, err := Parse(testdataDir("nested"))
	if err != nil {
		t.Fatal(err)
	}

	contactInfo := findStruct(t, structs, "ContactInfo")
	phone := findField(t, contactInfo, "Phone")
	if !phone.HasSanitize {
		t.Error("ContactInfo.Phone should have Sanitize()")
	}
	email := findField(t, contactInfo, "Email")
	if !email.HasSanitize {
		t.Error("ContactInfo.Email should have Sanitize()")
	}

	user := findStruct(t, structs, "User")
	ci := findField(t, user, "ContactInfo")
	if !ci.IsStruct {
		t.Error("User.ContactInfo should be a struct")
	}

	addr := findField(t, user, "Address")
	if !addr.IsStruct {
		t.Error("User.Address should be a struct")
	}
}

func TestParse_PointerAndSliceFields(t *testing.T) {
	_, structs, err := Parse(testdataDir("pointer_slice"))
	if err != nil {
		t.Fatal(err)
	}

	user := findStruct(t, structs, "User")

	phones := findField(t, user, "Phones")
	if !phones.IsSlice {
		t.Error("Phones should be a slice")
	}
	if !phones.HasSanitize {
		t.Error("Phones element should have Sanitize()")
	}

	address := findField(t, user, "Address")
	if !address.IsPtr {
		t.Error("Address should be a pointer")
	}
	if !address.IsStruct {
		t.Error("Address should be a struct")
	}
}

func TestParse_NoMaskableFields(t *testing.T) {
	_, structs, err := Parse(testdataDir("no_mask"))
	if err != nil {
		t.Fatal(err)
	}

	config := findStruct(t, structs, "Config")
	for _, f := range config.Fields {
		if f.HasSanitize {
			t.Errorf("field %s should not have Sanitize()", f.Name)
		}
	}
}

func TestParse_SkipsNonStructTypes(t *testing.T) {
	_, structs, err := Parse(testdataDir("basic"))
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range structs {
		if s.Name == "Phone" {
			t.Error("Phone is a named type, not a struct — should not be in struct list")
		}
	}
}

func findStruct(t *testing.T, structs []StructInfo, name string) StructInfo {
	t.Helper()
	for _, s := range structs {
		if s.Name == name {
			return s
		}
	}
	t.Fatalf("struct %q not found", name)
	return StructInfo{}
}

func findField(t *testing.T, s StructInfo, name string) FieldInfo {
	t.Helper()
	for _, f := range s.Fields {
		if f.Name == name {
			return f
		}
	}
	t.Fatalf("field %q not found in struct %q", name, s.Name)
	return FieldInfo{}
}
