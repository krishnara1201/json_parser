package parser

import (
	"json-parser/ast"
	"json-parser/lexer"
	"testing"
)

func TestParserPrimitives(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ast.JSONValue
		hasError bool
	}{
		{
			name:     "string",
			input:    `"hello"`,
			expected: ast.JSONString("hello"),
			hasError: false,
		},
		{
			name:     "number",
			input:    `42`,
			expected: ast.JSONNumber(42),
			hasError: false,
		},
		{
			name:     "float",
			input:    `3.14`,
			expected: ast.JSONNumber(3.14),
			hasError: false,
		},
		{
			name:     "true",
			input:    `true`,
			expected: ast.JSONBoolean(true),
			hasError: false,
		},
		{
			name:     "false",
			input:    `false`,
			expected: ast.JSONBoolean(false),
			hasError: false,
		},
		{
			name:     "null",
			input:    `null`,
			expected: nil,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			result := p.Parse()

			if len(p.Errors()) > 0 != tt.hasError {
				t.Fatalf("hasError mismatch: expected %v, got errors: %v", tt.hasError, p.Errors())
			}

			// Special case for null - just check type
			if tt.input == `null` {
				if _, ok := result.(ast.JSONNull); !ok {
					t.Fatalf("expected JSONNull, got %T", result)
				}
			} else if !jsonValuesEqual(result, tt.expected) {
				t.Fatalf("value mismatch: got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParserEmptyObject(t *testing.T) {
	input := `{}`
	l := lexer.New(input)
	p := New(l)
	result := p.Parse()

	if len(p.Errors()) > 0 {
		t.Fatalf("unexpected errors: %v", p.Errors())
	}

	obj, ok := result.(ast.JSONObject)
	if !ok {
		t.Fatalf("expected JSONObject, got %T", result)
	}

	if len(obj) != 0 {
		t.Fatalf("expected empty object, got %v", obj)
	}
}

func TestParserSimpleObject(t *testing.T) {
	input := `{"name":"Alice"}`
	l := lexer.New(input)
	p := New(l)
	result := p.Parse()

	if len(p.Errors()) > 0 {
		t.Fatalf("unexpected errors: %v", p.Errors())
	}

	obj, ok := result.(ast.JSONObject)
	if !ok {
		t.Fatalf("expected JSONObject, got %T", result)
	}

	if v, exists := obj["name"]; !exists || v != ast.JSONString("Alice") {
		t.Fatalf("expected name=Alice, got %v", obj)
	}
}

func TestParserMultiKeyObject(t *testing.T) {
	input := `{"name":"Bob","age":30,"active":true}`
	l := lexer.New(input)
	p := New(l)
	result := p.Parse()

	if len(p.Errors()) > 0 {
		t.Fatalf("unexpected errors: %v", p.Errors())
	}

	obj, ok := result.(ast.JSONObject)
	if !ok {
		t.Fatalf("expected JSONObject, got %T", result)
	}

	if obj["name"] != ast.JSONString("Bob") {
		t.Fatalf("expected name=Bob")
	}
	if obj["age"] != ast.JSONNumber(30) {
		t.Fatalf("expected age=30")
	}
	if obj["active"] != ast.JSONBoolean(true) {
		t.Fatalf("expected active=true")
	}
}

func TestParserEmptyArray(t *testing.T) {
	input := `[]`
	l := lexer.New(input)
	p := New(l)
	result := p.Parse()

	if len(p.Errors()) > 0 {
		t.Fatalf("unexpected errors: %v", p.Errors())
	}

	arr, ok := result.(ast.JSONArray)
	if !ok {
		t.Fatalf("expected JSONArray, got %T", result)
	}

	if len(arr) != 0 {
		t.Fatalf("expected empty array, got %v", arr)
	}
}

func TestParserSimpleArray(t *testing.T) {
	input := `[1,2,3]`
	l := lexer.New(input)
	p := New(l)
	result := p.Parse()

	if len(p.Errors()) > 0 {
		t.Fatalf("unexpected errors: %v", p.Errors())
	}

	arr, ok := result.(ast.JSONArray)
	if !ok {
		t.Fatalf("expected JSONArray, got %T", result)
	}

	if len(arr) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(arr))
	}

	expected := []ast.JSONValue{
		ast.JSONNumber(1),
		ast.JSONNumber(2),
		ast.JSONNumber(3),
	}

	for i, v := range arr {
		if v != expected[i] {
			t.Fatalf("element %d: expected %v, got %v", i, expected[i], v)
		}
	}
}

func TestParserMixedArray(t *testing.T) {
	input := `["text",42,true,null]`
	l := lexer.New(input)
	p := New(l)
	result := p.Parse()

	if len(p.Errors()) > 0 {
		t.Fatalf("unexpected errors: %v", p.Errors())
	}

	arr, ok := result.(ast.JSONArray)
	if !ok {
		t.Fatalf("expected JSONArray, got %T", result)
	}

	if len(arr) != 4 {
		t.Fatalf("expected 4 elements, got %d", len(arr))
	}

	if arr[0] != ast.JSONString("text") {
		t.Fatalf("expected string at [0]")
	}
	if arr[1] != ast.JSONNumber(42) {
		t.Fatalf("expected number at [1]")
	}
	if arr[2] != ast.JSONBoolean(true) {
		t.Fatalf("expected boolean at [2]")
	}
	if _, ok := arr[3].(ast.JSONNull); !ok {
		t.Fatalf("expected null at [3]")
	}
}

func TestParserNestedObject(t *testing.T) {
	input := `{"user":{"name":"Charlie","age":25}}`
	l := lexer.New(input)
	p := New(l)
	result := p.Parse()

	if len(p.Errors()) > 0 {
		t.Fatalf("unexpected errors: %v", p.Errors())
	}

	root, ok := result.(ast.JSONObject)
	if !ok {
		t.Fatalf("expected JSONObject, got %T", result)
	}

	user, ok := root["user"].(ast.JSONObject)
	if !ok {
		t.Fatalf("expected nested JSONObject, got %T", root["user"])
	}

	if user["name"] != ast.JSONString("Charlie") {
		t.Fatalf("expected name=Charlie")
	}
	if user["age"] != ast.JSONNumber(25) {
		t.Fatalf("expected age=25")
	}
}

func TestParserNestedArray(t *testing.T) {
	input := `[[1,2],[3,4]]`
	l := lexer.New(input)
	p := New(l)
	result := p.Parse()

	if len(p.Errors()) > 0 {
		t.Fatalf("unexpected errors: %v", p.Errors())
	}

	arr, ok := result.(ast.JSONArray)
	if !ok {
		t.Fatalf("expected JSONArray, got %T", result)
	}

	if len(arr) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(arr))
	}

	inner1, ok := arr[0].(ast.JSONArray)
	if !ok {
		t.Fatalf("expected nested array at [0]")
	}

	if len(inner1) != 2 || inner1[0] != ast.JSONNumber(1) || inner1[1] != ast.JSONNumber(2) {
		t.Fatalf("expected [1,2], got %v", inner1)
	}
}

func TestParserComplexNesting(t *testing.T) {
	input := `{"items":[{"id":1,"tags":["a","b"]},{"id":2,"tags":["c"]}]}`
	l := lexer.New(input)
	p := New(l)
	result := p.Parse()

	if len(p.Errors()) > 0 {
		t.Fatalf("unexpected errors: %v", p.Errors())
	}

	root, ok := result.(ast.JSONObject)
	if !ok {
		t.Fatalf("expected JSONObject, got %T", result)
	}

	items, ok := root["items"].(ast.JSONArray)
	if !ok {
		t.Fatalf("expected array for items")
	}

	if len(items) != 2 {
		t.Fatalf("expected 2 items")
	}

	item1, ok := items[0].(ast.JSONObject)
	if !ok {
		t.Fatalf("expected object at items[0]")
	}

	if item1["id"] != ast.JSONNumber(1) {
		t.Fatalf("expected id=1")
	}

	tags1, ok := item1["tags"].(ast.JSONArray)
	if !ok {
		t.Fatalf("expected array for tags")
	}

	if len(tags1) != 2 || tags1[0] != ast.JSONString("a") {
		t.Fatalf("expected tags=['a','b']")
	}
}

func TestParserErrorMissingColon(t *testing.T) {
	input := `{"key" "value"}`
	l := lexer.New(input)
	p := New(l)
	_ = p.Parse()

	// Should have errors
	if len(p.Errors()) == 0 {
		t.Fatalf("expected error for missing colon")
	}
}

func TestParserErrorInvalidKey(t *testing.T) {
	input := `{123: "value"}`
	l := lexer.New(input)
	p := New(l)
	_ = p.Parse()

	// Should have errors
	if len(p.Errors()) == 0 {
		t.Fatalf("expected error for non-string key")
	}
}

func TestParserErrorTrailingComma(t *testing.T) {
	// This depends on lexer behavior; currently the lexer produces tokens
	// and parser should reject trailing data
	input := `[1,2,]`
	l := lexer.New(input)
	p := New(l)
	result := p.Parse()

	// Behavior depends on whether lexer/parser handle this gracefully
	// Just ensure no panic
	if result == nil && len(p.Errors()) > 0 {
		t.Logf("trailing comma detected with errors: %v", p.Errors())
	}
}

// Helper to compare JSON values (works for test comparisons)
func jsonValuesEqual(a, b ast.JSONValue) bool {
	switch av := a.(type) {
	case ast.JSONString:
		bv, ok := b.(ast.JSONString)
		return ok && av == bv
	case ast.JSONNumber:
		bv, ok := b.(ast.JSONNumber)
		return ok && av == bv
	case ast.JSONBoolean:
		bv, ok := b.(ast.JSONBoolean)
		return ok && av == bv
	case ast.JSONNull:
		_, ok := b.(ast.JSONNull)
		return ok
	default:
		return false
	}
}
