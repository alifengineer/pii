package pii

import (
	"encoding/json"
	"reflect"
	"testing"
)

type UserPhone string

type User struct {
	Name        *string           `json:"name" pii:"name"`
	Phone       *UserPhone        `json:"phone" pii:"phone"`
	Email       string            `json:"email" pii:"email"`
	ContactInfo ContactInfo       `json:"contact_info"`
	Address     *Address          `json:"address"`
	UserID      int               `json:"user_id" pii:"ID"`
	ExtraInfo   map[string]string `json:"extra_info" pii:"email"`
}

type Address struct {
	Street string `json:"street"`
	Zip    string `json:"zip"`
}

type ContactInfo struct {
	Phone string `json:"phone" pii:"phone"`
	Email string `json:"email" pii:"email"`
}

func TestMask(t *testing.T) {
	type testCase struct {
		name     string
		input    any
		expected any
	}

	Register(TagValue("ID"), func(s string) string {
		return "****"
	})

	testCases := []testCase{
		{
			name: "masking struct with pii tags",
			input: User{
				UserID: 10,
				Name:   new("John Doe"),
				Phone:  new(UserPhone("+1234567890")),
				Email:  "john.doe@example.com",
				ContactInfo: ContactInfo{
					Phone: "+1234567890",
					Email: "jane.doe@example.com",
				},
				Address: &Address{
					Street: "123 Main St",
					Zip:    "12345",
				},
			},
			expected: User{
				UserID: 10,
				Name:   new("J*** D**"),
				Phone:  new(UserPhone("***7890")),
				Email:  "j***@e***.com",
				ContactInfo: ContactInfo{
					Phone: "***7890",
					Email: "j***@e***.com",
				},
				Address: &Address{
					Street: "123 Main St",
					Zip:    "12345",
				},
			},
		},
		{
			name: "nil pointer fields stay nil",
			input: User{
				Name:  nil,
				Phone: nil,
				Email: "a@b.com",
			},
			expected: User{
				Name:  nil,
				Phone: nil,
				Email: "a***@b***.com",
			},
		},
		{
			name: "all fields empty strings",
			input: User{
				Name:  new(""),
				Phone: new(UserPhone("")),
				Email: "",
				ContactInfo: ContactInfo{
					Phone: "",
					Email: "",
				},
			},
			expected: User{
				Name:  new(""),
				Phone: new(UserPhone("")),
				Email: "",
				ContactInfo: ContactInfo{
					Phone: "",
					Email: "",
				},
			},
		},
		{
			name: "single character values",
			input: User{
				Name:  new("A"),
				Phone: new(UserPhone("5")),
				Email: "x",
				ContactInfo: ContactInfo{
					Phone: "5",
					Email: "x",
				},
			},
			expected: User{
				Name:  new("A"),
				Phone: new(UserPhone("****")),
				Email: "****",
				ContactInfo: ContactInfo{
					Phone: "****",
					Email: "****",
				},
			},
		},
		{
			name: "pointer to struct passed to Mask",
			input: &User{
				Name:  new("John"),
				Phone: new(UserPhone("+1234567890")),
				Email: "j@d.com",
			},
			expected: User{
				Name:  new("J***"),
				Phone: new(UserPhone("***7890")),
				Email: "j***@d***.com",
			},
		},
		{
			name: "struct with no pii tags unchanged",
			input: Address{
				Street: "123 Main St",
				Zip:    "12345",
			},
			expected: Address{
				Street: "123 Main St",
				Zip:    "12345",
			},
		},
		{
			name: "nil pointer to nested struct",
			input: User{
				Email:   "a@b.com",
				Address: nil,
			},
			expected: User{
				Email:   "a***@b***.com",
				Address: nil,
			},
		},
		{
			name: "unicode names",
			input: User{
				Name:  new("Алиса Иванова"),
				Email: "alice@example.com",
			},
			expected: User{
				Name:  new("А**** И******"),
				Email: "a***@e***.com",
			},
		},
		{
			name: "non-string pii tagged field (int) unchanged",
			input: User{
				UserID: 42,
				Email:  "a@b.com",
			},
			expected: User{
				UserID: 42,
				Email:  "a***@b***.com",
			},
		},
		{
			name:     "non-struct input returns as-is",
			input:    "just a string",
			expected: "just a string",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Mask(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				jsonInput, _ := json.MarshalIndent(tc.input, "", "  ")
				t.Logf("input:\n%s", jsonInput)

				jsonExpected, _ := json.MarshalIndent(tc.expected, "", "  ")
				t.Logf("expected:\n%s", jsonExpected)

				jsonResult, _ := json.MarshalIndent(result, "", "  ")
				t.Logf("got:\n%s", jsonResult)

				t.Errorf("expected %+v, got %+v", tc.expected, result)
			}

			t.Log(tc.input)
		})
	}
}

var benchUser = User{
	Name:  new("John Doe"),
	Phone: new(UserPhone("+1234567890")),
	Email: "john.doe@example.com",
	ContactInfo: ContactInfo{
		Phone: "+1234567890",
		Email: "jane.doe@example.com",
	},
	Address: &Address{
		Street: "123 Main St",
		Zip:    "12345",
	},
}

func BenchmarkJSONMarshal(b *testing.B) {
	for b.Loop() {
		json.Marshal(benchUser) //nolint:errcheck,gosec //best effort
	}
}

func BenchmarkMask(b *testing.B) {
	for b.Loop() {
		Mask(benchUser)
	}
}

func BenchmarkMaskThenJSON(b *testing.B) {
	for b.Loop() {
		json.Marshal(Mask(benchUser)) //nolint:errcheck,gosec //best effort
	}
}

func BenchmarkMaskParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Mask(benchUser)
		}
	})
}

func BenchmarkMaskName(b *testing.B) {
	for b.Loop() {
		MaskName("Jonathan")
	}
}

func BenchmarkMaskPhone(b *testing.B) {
	for b.Loop() {
		MaskPhone("+998901234567")
	}
}

func BenchmarkMaskEmail(b *testing.B) {
	for b.Loop() {
		MaskEmail("john.doe@example.com")
	}
}

func BenchmarkMaskPAN(b *testing.B) {
	for b.Loop() {
		MaskPAN("4111111111111111")
	}
}
