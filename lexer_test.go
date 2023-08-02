package gogqllexer

import (
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

func TestLexer_NextToken_ReadSingleName(t *testing.T) {
	tests := []struct {
		name string
		src  *Source
		want []Token
	}{
		{
			name: "simple name",
			src: &Source{
				Body: "query",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Name,
					Value: "query",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "simple name",
			src: &Source{
				Body: "_query",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Name,
					Value: "_query",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "simple name",
			src: &Source{
				Body: "_0query",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Name,
					Value: "_0query",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "white space",
			src: &Source{
				Body: "  query  ",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Name,
					Value: "query",
					Position: Position{
						Line:  1,
						Start: 3,
					},
				},
			},
		},
		{
			name: "line feed",
			src: &Source{
				Body: "\nquery",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Name,
					Value: "query",
					Position: Position{
						Line:  2,
						Start: 2,
					},
				},
			},
		},
		{
			name: "carriage return",
			src: &Source{
				Body: "\rquery",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Name,
					Value: "query",
					Position: Position{
						Line:  2,
						Start: 2,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.src, strings.NewReader(tt.src.Body))

			gotTokens := make([]Token, 0)
			for {
				got, err := l.NextToken()
				if err != nil {
					t.Fatal(err)
				}
				if got.Kind == EOF {
					t.Log(got)
					break
				}

				gotTokens = append(gotTokens, got)
			}

			ok := assert.Equal(t, tt.want, gotTokens)
			if !ok {
				t.Fatal("miss")
			}
		})
	}
}

func TestLexer_NextToken_SinglePunctuator(t *testing.T) {
	tests := []struct {
		name string
		src  *Source
		want []Token
	}{
		{
			name: "punctuator bang",
			src: &Source{
				Body: "!",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Bang,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "punctuator dollar",
			src: &Source{
				Body: "$",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Dollar,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "punctuator amp",
			src: &Source{
				Body: "&",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Amp,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "punctuator paren left",
			src: &Source{
				Body: "(",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  ParenL,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "punctuator paren right",
			src: &Source{
				Body: ")",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  ParenR,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "punctuator spread",
			src: &Source{
				Body: "...",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Spread,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "punctuator colon",
			src: &Source{
				Body: ":",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Colon,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "punctuator equal",
			src: &Source{
				Body: "=",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Equal,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "punctuator at",
			src: &Source{
				Body: "@",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  At,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "punctuator bracket left",
			src: &Source{
				Body: "[",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  BracketL,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "punctuator bracket right",
			src: &Source{
				Body: "]",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  BracketR,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "punctuator brace left",
			src: &Source{
				Body: "{",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  BraceL,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "punctuator brace right",
			src: &Source{
				Body: "}",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  BraceR,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "punctuator pipe",
			src: &Source{
				Body: "|",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Pipe,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.src, strings.NewReader(tt.src.Body))

			gotTokens := make([]Token, 0)
			for {
				got, err := l.NextToken()
				if err != nil {
					t.Fatal(err)
				}
				if got.Kind == EOF {
					break
				}

				gotTokens = append(gotTokens, got)
			}

			ok := assert.Equal(t, tt.want, gotTokens)
			if !ok {
				t.Fatal("miss")
			}
		})
	}
}

func TestLexer_NextToken_Comment(t *testing.T) {
	tests := []struct {
		name string
		src  *Source
		want []Token
	}{
		{
			name: "read comment token",
			src: &Source{
				Body: "# This is comment.",
				Name: "Spec_IgnoreWhiteSpace",
			},
			want: []Token{
				{
					Kind:  Comment,
					Value: "# This is comment.",
					Position: Position{
						Line:  1,
						Start: 0,
					},
				},
			},
		},
		{
			name: "read comment token",
			src: &Source{
				Body: "# This is comment.\n\r\n",
				Name: "Spec_IgnoreWhiteSpace",
			},
			want: []Token{
				{
					Kind:  Comment,
					Value: "# This is comment.",
					Position: Position{
						Line:  1,
						Start: 0,
					},
				},
			},
		},
		{
			name: "read comment token",
			src: &Source{
				Body: "\n\r\n# This is comment.",
				Name: "Spec_IgnoreWhiteSpace",
			},
			want: []Token{
				{
					Kind:  Comment,
					Value: "# This is comment.",
					Position: Position{
						Line:  3,
						Start: 3,
					},
				},
			},
		},
		{
			name: "read comment token",
			src: &Source{
				Body: "# This is comment.   ",
				Name: "Spec_IgnoreWhiteSpace",
			},
			want: []Token{
				{
					Kind:  Comment,
					Value: "# This is comment.   ",
					Position: Position{
						Line:  1,
						Start: 0,
					},
				},
			},
		},
		{
			name: "read comment token",
			src: &Source{
				Body: "# This is first comment.\n# This is second comment.",
				Name: "Spec_IgnoreWhiteSpace",
			},
			want: []Token{
				{
					Kind:  Comment,
					Value: "# This is first comment.",
					Position: Position{
						Line:  1,
						Start: 0,
					},
				},
				{
					Kind:  Comment,
					Value: "# This is second comment.",
					Position: Position{
						Line:  2,
						Start: 25,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				src:  tt.src,
				line: 1,
			}

			gotTokens := make([]Token, 0)
			for {
				got, err := l.NextToken()
				if err != nil {
					t.Fatal(err)
				}
				if got.Kind == EOF {
					t.Log(got)
					break
				}

				gotTokens = append(gotTokens, got)
			}

			ok := assert.Equal(t, tt.want, gotTokens)
			if !ok {
				t.Fatal("miss")
			}
		})
	}
}

func TestLexer_NextToken_Int(t *testing.T) {
	tests := []struct {
		name string
		src  *Source
		want []Token
	}{
		{
			name: "IntToken_0",
			src: &Source{
				Body: "0",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Int,
					Value: "0",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "IntToken_1",
			src: &Source{
				Body: "1",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Int,
					Value: "1",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "IntToken_9",
			src: &Source{
				Body: "9",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Int,
					Value: "9",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "IntToken_100",
			src: &Source{
				Body: "100",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Int,
					Value: "100",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "IntToken_Negative",
			src: &Source{
				Body: "-9",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Int,
					Value: "-9",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.src, strings.NewReader(tt.src.Body))

			gotTokens := make([]Token, 0)
			for {
				got, err := l.NextToken()
				if err != nil {
					t.Fatal(err)
				}
				if got.Kind == EOF {
					break
				}

				gotTokens = append(gotTokens, got)
			}

			ok := assert.Equal(t, tt.want, gotTokens)
			if !ok {
				t.Fatal("miss")
			}
		})
	}
}

func TestLexer_NextToken_Float(t *testing.T) {
	tests := []struct {
		name string
		src  *Source
		want []Token
	}{
		{
			name: "FloatToken",
			src: &Source{
				Body: "0.1",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Float,
					Value: "0.1",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "FloatToken",
			src: &Source{
				Body: "0.100",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Float,
					Value: "0.100",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "FloatToken",
			src: &Source{
				Body: "0.0021",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Float,
					Value: "0.0021",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "FloatToken",
			src: &Source{
				Body: "123.0021",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Float,
					Value: "123.0021",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "FloatToken",
			src: &Source{
				Body: "-123.0021",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Float,
					Value: "-123.0021",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "FloatToken",
			src: &Source{
				Body: "0.0",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Float,
					Value: "0.0",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.src, strings.NewReader(tt.src.Body))

			gotTokens := make([]Token, 0)
			for {
				got, err := l.NextToken()
				if err != nil {
					t.Fatal(err)
				}
				if got.Kind == EOF {
					t.Log(got)
					break
				}

				gotTokens = append(gotTokens, got)
			}

			ok := assert.Equal(t, tt.want, gotTokens)
			if !ok {
				t.Fatal("miss")
			}
		})
	}
}

func TestLexer_NextToken_Exponent(t *testing.T) {
	tests := []struct {
		name string
		src  *Source
		want []Token
	}{
		{
			name: "ExponentToken",
			src: &Source{
				Body: "1e50",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Float,
					Value: "1e50",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "ExponentToken",
			src: &Source{
				Body: "1.0e50",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Float,
					Value: "1.0e50",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "ExponentToken",
			src: &Source{
				Body: "1.0e-50",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Float,
					Value: "1.0e-50",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "ExponentToken",
			src: &Source{
				Body: "1.0e+50",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Float,
					Value: "1.0e+50",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.src, strings.NewReader(tt.src.Body))

			gotTokens := make([]Token, 0)
			for {
				got, err := l.NextToken()
				if err != nil {
					t.Fatal(err)
				}
				if got.Kind == EOF {
					t.Log(got)
					break
				}

				gotTokens = append(gotTokens, got)
			}

			ok := assert.Equal(t, tt.want, gotTokens)
			if !ok {
				t.Fatal("miss")
			}
		})
	}
}

func TestLexer_NextToken_String_Invalid(t *testing.T) {
	tests := []struct {
		name string
		src  *Source
		want []Token
	}{
		{
			name: "not closing string value",
			src: &Source{
				Body: "\"not closing string value",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Invalid,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "invalid string character (line feed)",
			src: &Source{
				Body: "\"\n\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Invalid,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "invalid string character (line carriage return)",
			src: &Source{
				Body: "\"\r\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Invalid,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "invalid string character (single backslash)",
			src: &Source{
				Body: "\"\\\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Invalid,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		// invalid escaped character
		{
			name: "escaped character (backslash)",
			src: &Source{
				Body: "\"\\\\\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\\\\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		// invalid escaped unicode
		{
			name: "escaped unicode over f",
			src: &Source{
				Body: "\"\\u000g\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Invalid,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "escaped unicode less than 4 digits",
			src: &Source{
				Body: "\"\\u000\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Invalid,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "escaped unicode less than 4 digits",
			src: &Source{
				Body: "\"\\u\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  Invalid,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				src:  tt.src,
				line: 1,
			}

			gotTokens := make([]Token, 0)
			for {
				got, err := l.NextToken()
				if err != nil {
					t.Fatal(err)
				}
				if got.Kind == EOF {
					break
				}

				gotTokens = append(gotTokens, got)
				if got.Kind == Invalid {
					break
				}
			}

			ok := assert.Equal(t, tt.want, gotTokens)
			if !ok {
				t.Fatal("miss")
			}
		})
	}
}

func TestLexer_NextToken_String(t *testing.T) {
	tests := []struct {
		name string
		src  *Source
		want []Token
	}{
		{
			name: "empty string",
			src: &Source{
				Body: "\"\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  String,
					Value: "\"\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "simple string",
			src: &Source{
				Body: "\"simple string\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  String,
					Value: "\"simple string\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "simple string with white space",
			src: &Source{
				Body: "\"  simple string  \"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  String,
					Value: "\"  simple string  \"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		// escaped character
		{
			name: "escaped character (backslash)",
			src: &Source{
				Body: "\"\\\\\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\\\\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "escaped character (double quote)",
			src: &Source{
				Body: "\"\\\"\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\\"\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "escaped character (slash)",
			src: &Source{
				Body: "\"\\/\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\/\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "escaped character (backspace)",
			src: &Source{
				Body: "\"\\b\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\b\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "escaped character (form feed)",
			src: &Source{
				Body: "\"\\f\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\f\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "escaped character (line feed)",
			src: &Source{
				Body: "\"\\n\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\n\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "escaped character (carriage return)",
			src: &Source{
				Body: "\"\\r\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\r\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "escaped character (horizontal tab)",
			src: &Source{
				Body: "\"\\t\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\t\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		// escaped unicode
		{
			name: "escaped unicode",
			src: &Source{
				Body: "\"\\u000a\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\u000a\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "escaped unicode",
			src: &Source{
				Body: "\"\\u0000\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\u0000\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "escaped unicode",
			src: &Source{
				Body: "\"\\uffff\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\uffff\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "escaped unicode",
			src: &Source{
				Body: "\"\\uffff0\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\uffff0\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				src:  tt.src,
				line: 1,
			}

			gotTokens := make([]Token, 0)
			for {
				got, err := l.NextToken()
				if err != nil {
					t.Fatal(err)
				}
				if got.Kind == EOF {
					break
				}

				gotTokens = append(gotTokens, got)
			}

			ok := assert.Equal(t, tt.want, gotTokens)
			if !ok {
				t.Fatal("miss")
			}
		})
	}
}

func TestLexer_NextToken_BlockString(t *testing.T) {
	tests := []struct {
		name string
		src  *Source
		want []Token
	}{
		{
			name: "empty block string",
			src: &Source{
				Body: "\"\"\"\"\"\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  BlockString,
					Value: "\"\"\"\"\"\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "simple string",
			src: &Source{
				Body: "\"\"\"simple string\"\"\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  BlockString,
					Value: "\"\"\"simple string\"\"\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "white space",
			src: &Source{
				Body: "\"\"\"  simple string  \"\"\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  BlockString,
					Value: "\"\"\"  simple string  \"\"\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "line feed",
			src: &Source{
				Body: "\"\"\" \nsimple string\"\"\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  BlockString,
					Value: "\"\"\" \nsimple string\"\"\"",
					Position: Position{
						Line:  2,
						Start: 1,
					},
				},
			},
		},
		{
			name: "line carriage return",
			src: &Source{
				Body: "\"\"\" \rsimple string\"\"\"",
				Name: "Spec",
			},
			want: []Token{
				{
					Kind:  BlockString,
					Value: "\"\"\" \rsimple string\"\"\"",
					Position: Position{
						Line:  2,
						Start: 1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				src:  tt.src,
				line: 1,
			}

			gotTokens := make([]Token, 0)
			for {
				got, err := l.NextToken()
				if err != nil {
					t.Fatal(err)
				}
				if got.Kind == EOF {
					break
				}

				gotTokens = append(gotTokens, got)
			}

			ok := assert.Equal(t, tt.want, gotTokens)
			if !ok {
				t.Fatal("miss")
			}
		})
	}
}

func TestLexer(t *testing.T) {
	// testdata/schema/schema.graphqlを読み取る
	f, err := os.Open("./testdata/schema/schema.graphqls")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	src := &Source{
		Body: string(b),
		Name: "testdata/schema/schema.graphqls",
	}

	l := New(src, strings.NewReader(string(b)))

	tokens := make([]Token, 0)
	for {
		got, err := l.NextToken()
		if err != nil {
			t.Fatal(err)
		}

		tokens = append(tokens, got)
		if got.Kind == EOF || got.Kind == Invalid {
			break
		}
	}

	for _, token := range tokens {
		log.Println(DebugTokenString(token), token.Value, token.Position)
	}
}

// DebugTokenString はトークンの文字列表現を返す
func DebugTokenString(t Token) string {
	switch t.Kind {
	case EOF:
		return "EOF"
	case Invalid:
		return "Invalid"
	case Comment:
		return "Comment"
	case Int:
		return "Int"
	case Float:
		return "Float"
	case String:
		return "String"
	case BlockString:
		return "BlockString"
	case Name:
		return "Name"
	case Bang:
		return "!"
	case Dollar:
		return "$"
	case Amp:
		return "&"
	case ParenL:
		return "("
	case ParenR:
		return ")"
	case Spread:
		return "..."
	case Colon:
		return ":"
	case Equal:
		return "="
	case At:
		return "@"
	case BracketR:
		return "["
	case BracketL:
		return "]"
	case Pipe:
		return "|"
	case BraceL:
		return "{"
	case BraceR:
		return "}"
	default:
		return ""
	}
}
