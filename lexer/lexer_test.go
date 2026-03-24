package lexer

import "testing"

func TestNextTokenTable(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "punctuation only",
			input: "{}[]:,",
			expected: []Token{
				{Type: LBRACE, Value: "{"},
				{Type: RBRACE, Value: "}"},
				{Type: LBRACKET, Value: "["},
				{Type: RBRACKET, Value: "]"},
				{Type: COLON, Value: ":"},
				{Type: COMMA, Value: ","},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "string and number",
			input: "\"name\" 123",
			expected: []Token{
				{Type: STRING, Value: "name"},
				{Type: NUMBER, Value: "123"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "keywords",
			input: "true false null",
			expected: []Token{
				{Type: TRUE, Value: "true"},
				{Type: FALSE, Value: "false"},
				{Type: NULL, Value: "null"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name:  "mini object",
			input: "{\"k\":1}",
			expected: []Token{
				{Type: LBRACE, Value: "{"},
				{Type: STRING, Value: "k"},
				{Type: COLON, Value: ":"},
				{Type: NUMBER, Value: "1"},
				{Type: RBRACE, Value: "}"},
				{Type: EOF, Value: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{input: tt.input}
			l.readChar()

			for i, want := range tt.expected {
				got := l.NextToken()
				if got.Type != want.Type || got.Value != want.Value {
					t.Fatalf("token %d mismatch: got=(%s,%q) want=(%s,%q)", i, got.Type, got.Value, want.Type, want.Value)
				}
			}
		})
	}
}
