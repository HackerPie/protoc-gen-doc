package gendoc

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_generateSample(t *testing.T) {
	mapEntryMessage := &Message{
		FullName:  "a.b.c.MapEntry",
		HasFields: true,
		Fields: []*MessageField{
			{Name: "key", FullType: "string"},
			{Name: "value", FullType: "string"},
		},
	}

	baseMsg := &Message{
		FullName:  "a.b.c.Base",
		HasFields: true,
		Fields: []*MessageField{
			{Name: "double_field", FullType: "double"},
			{Name: "float_field", FullType: "float"},
			{Name: "map_field", IsMap: true, FullType: "a.b.c.MapEntry"},
		},
	}

	simpleMsg := &Message{
		FullName:  "a.b.c.Simple",
		HasFields: true,
		Fields: []*MessageField{
			{Name: "double_field", FullType: "double"},
		},
	}

	baseEnum := &Enum{
		FullName: "a.b.c.Enum",
		Values: []*EnumValue{
			{Name: "value1", Number: "1"},
			{Name: "value2", Number: "2"},
			{Name: "value3", Number: "3"},
		},
	}

	file := &File{
		Package:        "a.b.c",
		HasEnums:       true,
		HasMessages:    true,
		Enums:          orderedEnums{baseEnum},
		EnumIndexes:    map[string]int{"a.b.c.Enum": 0},
		Messages:       orderedMessages{mapEntryMessage, baseMsg, simpleMsg},
		MessageIndexes: map[string]int{"a.b.c.MapEntry": 0, "a.b.c.Base": 1, "a.b.c.Simple": 2},
	}

	message := &Message{
		FullName:  "a.b.c.Test",
		HasFields: true,
		HasOneofs: true,
		Fields: []*MessageField{
			{Name: "double_field", FullType: "double"},
			{Name: "float_field", FullType: "float"},
			{Name: "int64_field", FullType: "int64"},
			{Name: "uint32_field", FullType: "uint32"},
			{Name: "uint64_field", FullType: "uint64"},
			{Name: "sint32_field", FullType: "sint32"},
			{Name: "sint64_field", FullType: "sint64"},
			{Name: "fixed32_field", FullType: "fixed32"},
			{Name: "fixed64_field", FullType: "fixed64"},
			{Name: "sfixed32_field", FullType: "sfixed32"},
			{Name: "sfixed64_field", FullType: "sfixed64"},
			{Name: "bool_field", FullType: "bool", IsOneof: true, OneofDecl: "test_oneof"},
			{Name: "string_field", FullType: "string", IsOneof: true, OneofDecl: "test_oneof"},
			{Name: "bytes_field", FullType: "bytes"},
			{Name: "message_field", FullType: "a.b.c.Base"},
			{Name: "enum_field", FullType: "a.b.c.Enum"},
		},
	}

	got := generateSample(message, file)
	expect := map[string]interface{}{
		"double_field":   0.01,
		"float_field":    0.01,
		"int64_field":    "0",
		"uint32_field":   0,
		"uint64_field":   "0",
		"sint32_field":   0,
		"sint64_field":   "0",
		"fixed32_field":  0,
		"fixed64_field":  "0",
		"sfixed32_field": 0,
		"sfixed64_field": "0",
		"string_field":   "string",
		"bool_field":     false,
		//"string_field":   "string", // this field is omitted because of oneof
		"bytes_field": "<binary bytes>",
		"message_field": map[string]interface{}{
			"double_field": 0.01,
			"float_field":  0.01,
			"map_field": map[interface{}]interface{}{
				"string": "string",
			},
		},
		"enum_field": 3,
	}

	require.Equal(t, expect, got)
}
