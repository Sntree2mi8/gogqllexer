package gogqllexer

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestLexer_NextToken_SkipIgnoredTokens(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []Token
	}{
		{
			name: "ignore white space",
			src:  "  query  ",
			want: []Token{
				{
					Kind:  Name,
					Value: "query",
					Position: Position{
						Line:  1,
						Start: 3,
					},
				},
				{
					Kind:  EOF,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 9,
					},
				},
			},
		},
		{
			name: "ignore line terminator",
			src:  "\n\rquery\n\r\n",
			want: []Token{
				{
					Kind:  Name,
					Value: "query",
					Position: Position{
						Line:  3,
						Start: 3,
					},
				},
				{
					Kind:  EOF,
					Value: "",
					Position: Position{
						Line:  5,
						Start: 10,
					},
				},
			},
		},
		{
			name: "ignore unicode byte order mark",
			src:  "\uFEFFquery",
			want: []Token{
				{
					Kind:  Name,
					Value: "query",
					Position: Position{
						Line:  1,
						Start: 4,
					},
				},
				{
					Kind:  EOF,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 8,
					},
				},
			},
		},
		{
			name: "ignore comma",
			src:  ",query,",
			want: []Token{
				{
					Kind:  Name,
					Value: "query",
					Position: Position{
						Line:  1,
						Start: 2,
					},
				},
				{
					Kind:  EOF,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 7,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(strings.NewReader(tt.src))

			gotTokens := make([]Token, 0)
			for {
				got, err := l.NextToken()
				if err != nil {
					t.Fatal(err)
				}

				gotTokens = append(gotTokens, got)
				if got.Kind == EOF {
					break
				}
			}

			assert.Equal(t, tt.want, gotTokens)
		})
	}
}

func TestLexer_NextToken_Name(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []Token
	}{
		{
			name: "name",
			src:  "_queryQUERY0123456789$",
			want: []Token{
				{
					Kind:  Name,
					Value: "_queryQUERY0123456789",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind:  Dollar,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 22,
					},
				},
				{
					Kind:  EOF,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 22,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(strings.NewReader(tt.src))

			gotTokens := make([]Token, 0)
			for {
				got, err := l.NextToken()
				if err != nil {
					t.Fatal(err)
				}

				gotTokens = append(gotTokens, got)
				if got.Kind == EOF || got.Kind == Invalid {
					break
				}
			}

			assert.Equal(t, tt.want, gotTokens)
		})
	}
}

func TestLexer_NextToken_Punctuator(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []Token
	}{
		{
			name: "punctuator bang",
			src:  "!$&()...:=@[]{|}",
			want: []Token{
				{
					Kind:  Bang,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind:  Dollar,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 2,
					},
				},
				{
					Kind:  Amp,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 3,
					},
				},
				{
					Kind:  ParenL,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 4,
					},
				},
				{
					Kind:  ParenR,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 5,
					},
				},
				{
					Kind:  Spread,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 6,
					},
				},
				{
					Kind:  Colon,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 9,
					},
				},
				{
					Kind:  Equal,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 10,
					},
				},
				{
					Kind:  At,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 11,
					},
				},
				{
					Kind:  BracketL,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 12,
					},
				},
				{
					Kind:  BracketR,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 13,
					},
				},
				{
					Kind:  BraceL,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 14,
					},
				},
				{
					Kind:  Pipe,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 15,
					},
				},
				{
					Kind:  BraceR,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 16,
					},
				},
				{
					Kind:  EOF,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 16,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(strings.NewReader(tt.src))

			gotTokens := make([]Token, 0)
			for {
				got, err := l.NextToken()
				if err != nil {
					t.Fatal(err)
				}

				gotTokens = append(gotTokens, got)
				if got.Kind == EOF || got.Kind == Invalid {
					break
				}
			}

			assert.Equal(t, tt.want, gotTokens)
		})
	}
}

func TestLexer_NextToken_Comment(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []Token
	}{
		{
			name: "read comment token",
			src:  "# comment",
			want: []Token{
				{
					Kind:  Comment,
					Value: "# comment",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind:  EOF,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 9,
					},
				},
			},
		},
		{
			name: "read comment token",
			src:  "# comment\n\r\n",
			want: []Token{
				{
					Kind:  Comment,
					Value: "# comment",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind:  EOF,
					Value: "",
					Position: Position{
						Line:  3,
						Start: 12,
					},
				},
			},
		},
		{
			name: "read comment token",
			src:  "\n\r\n# comment",
			want: []Token{
				{
					Kind:  Comment,
					Value: "# comment",
					Position: Position{
						Line:  3,
						Start: 4,
					},
				},
				{
					Kind:  EOF,
					Value: "",
					Position: Position{
						Line:  3,
						Start: 12,
					},
				},
			},
		},
		{
			name: "read comment token",
			src:  "# comment   ",
			want: []Token{
				{
					Kind:  Comment,
					Value: "# comment   ",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind:  EOF,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 12,
					},
				},
			},
		},
		{
			name: "read comment token",
			src:  "# comment1 # comment1",
			want: []Token{
				{
					Kind:  Comment,
					Value: "# comment1 # comment1",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind:  EOF,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 21,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(strings.NewReader(tt.src))

			gotTokens := make([]Token, 0)
			for {
				got, err := l.NextToken()
				if err != nil {
					t.Fatal(err)
				}

				gotTokens = append(gotTokens, got)
				if got.Kind == EOF || got.Kind == Invalid {
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

func TestLexer_NextToken_ReadNumber(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []Token
	}{
		// int
		{
			name: "0",
			src:  "0",
			want: []Token{
				{
					Kind:  Int,
					Value: "0",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
			},
		},
		{
			name: "1",
			src:  "1",
			want: []Token{
				{
					Kind:  Int,
					Value: "1",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 1,
					},
				}},
		},
		{
			name: "9",
			src:  "9",
			want: []Token{
				{
					Kind:  Int,
					Value: "9",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 1,
					},
				}},
		},
		{
			name: "100",
			src:  "100",
			want: []Token{
				{
					Kind:  Int,
					Value: "100",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 3,
					},
				}},
		},
		{
			name: "negative",
			src:  "-9",
			want: []Token{
				{
					Kind:  Int,
					Value: "-9",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 2,
					},
				}},
		},
		// int invalid
		{
			name: "IntValue must not any leading 0",
			src:  "0123",
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
			name: "Int token can't end with dot",
			src:  "1.",
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
			name: "Int token can't end with dot",
			src:  "0.",
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
			name: "Int token can't end with name start character",
			src:  "1a",
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
			name: "Int token can't end with name start character",
			src:  "0a",
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
			name: "Int token can't end with name start character",
			src:  "1_",
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
			name: "Int token can't end with name start character",
			src:  "0_",
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
		// float(fractional)
		{
			name: "FloatToken",
			src:  "0.1",
			want: []Token{
				{
					Kind:  Float,
					Value: "0.1",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 3,
					},
				},
			},
		},
		{
			name: "FloatToken",
			src:  "0.100",
			want: []Token{
				{
					Kind:  Float,
					Value: "0.100",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 5,
					},
				},
			},
		},
		{
			name: "FloatToken",
			src:  "0.0021",
			want: []Token{
				{
					Kind:  Float,
					Value: "0.0021",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 6,
					},
				},
			},
		},
		{
			name: "FloatToken",
			src:  "123.0021",
			want: []Token{
				{
					Kind:  Float,
					Value: "123.0021",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 8,
					},
				},
			},
		},
		{
			name: "FloatToken",
			src:  "-123.0021",
			want: []Token{
				{
					Kind:  Float,
					Value: "-123.0021",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 9,
					},
				},
			},
		},
		{
			name: "FloatToken",
			src:  "0.0",
			want: []Token{
				{
					Kind:  Float,
					Value: "0.0",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 3,
					},
				},
			},
		},
		// float(fractional) invalid
		{
			name: "FloatToken can't end with dot",
			src:  "0.1.",
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
			name: "FloatToken can't end with name start character",
			src:  "0.1a",
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
			name: "FloatToken can't end with name start character",
			src:  "0.1_",
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
		// float(exponent)
		{
			name: "ExponentToken",
			src:  "1e50",
			want: []Token{
				{
					Kind:  Float,
					Value: "1e50",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 4,
					},
				},
			},
		},
		{
			name: "ExponentToken",
			src:  "1.0e50",
			want: []Token{
				{
					Kind:  Float,
					Value: "1.0e50",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 6,
					},
				},
			},
		},
		{
			name: "ExponentToken",
			src:  "1.0e-50",
			want: []Token{
				{
					Kind:  Float,
					Value: "1.0e-50",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 7,
					},
				},
			},
		},
		{
			name: "ExponentToken",
			src:  "1.0e+50",
			want: []Token{
				{
					Kind:  Float,
					Value: "1.0e+50",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 7,
					},
				},
			},
		},
		// float(exponent) invalid
		{
			name: "ExponentToken can't end with dot",
			src:  "1e50.",
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
			name: "ExponentToken can't end with name start character",
			src:  "1e50a",
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
			name: "ExponentToken can't end with name start character",
			src:  "1e50_",
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
			l := New(strings.NewReader(tt.src))

			gotTokens := make([]Token, 0)
			for {
				got, err := l.NextToken()
				if err != nil {
					t.Fatal(err)
				}

				gotTokens = append(gotTokens, got)
				if got.Kind == EOF || got.Kind == Invalid {
					break
				}
			}

			assert.Equal(t, tt.want, gotTokens)
		})
	}
}

func TestLexer_NextToken_String(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want []Token
	}{
		// string
		{
			name: "empty string",
			src:  "\"\"",
			want: []Token{
				{
					Kind:  String,
					Value: "\"\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 2,
					},
				},
			},
		},
		{
			name: "simple string",
			src:  "\"simple string\"",
			want: []Token{
				{
					Kind:  String,
					Value: "\"simple string\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 15,
					},
				},
			},
		},
		{
			name: "simple string with white space",
			src:  "\"  simple string  \"",
			want: []Token{
				{
					Kind:  String,
					Value: "\"  simple string  \"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 19,
					},
				},
			},
		},
		// string escaped character
		{
			name: "escaped character (backslash)",
			src:  "\"\\\\\"",
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\\\\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 4,
					},
				},
			},
		},
		{
			name: "escaped character (double quote)",
			src:  "\"\\\"\"",
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\\"\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 4,
					},
				},
			},
		},
		{
			name: "escaped character (slash)",
			src:  "\"\\/\"",
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\/\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 4,
					},
				},
			},
		},
		{
			name: "escaped character (backspace)",
			src:  "\"\\b\"",
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\b\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 4,
					},
				},
			},
		},
		{
			name: "escaped character (form feed)",
			src:  "\"\\f\"",
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\f\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 4,
					},
				},
			},
		},
		{
			name: "escaped character (line feed)",
			src:  "\"\\n\"",
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\n\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 4,
					},
				},
			},
		},
		{
			name: "escaped character (carriage return)",
			src:  "\"\\r\"",
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\r\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 4,
					},
				},
			},
		},
		{
			name: "escaped character (horizontal tab)",
			src:  "\"\\t\"",
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\t\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 4,
					},
				},
			},
		},
		// string escaped unicode
		{
			name: "escaped unicode",
			src:  "\"\\u000a\"",
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\u000a\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 8,
					},
				},
			},
		},
		{
			name: "escaped unicode",
			src:  "\"\\u0000\"",
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\u0000\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 8,
					},
				},
			},
		},
		{
			name: "escaped unicode",
			src:  "\"\\uffff\"",
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\uffff\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 8,
					},
				},
			},
		},
		{
			name: "escaped unicode",
			src:  "\"\\uffff0\"",
			want: []Token{
				{
					Kind:  String,
					Value: "\"\\uffff0\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 9,
					},
				},
			},
		},
		// string invalid
		{
			name: "not closing string value",
			src:  "\"not closing string value",
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
			src:  "\"\n\"",
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
			src:  "\"\r\"",
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
			src:  "\"\\\"",
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
		// string invalid escaped unicode
		{
			name: "escaped unicode over f",
			src:  "\"\\u000g\"",
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
			src:  "\"\\u000\"",
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
			src:  "\"\\u\"",
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
		// block string
		{
			name: "empty block string",
			src:  "\"\"\"\"\"\"",
			want: []Token{
				{
					Kind:  BlockString,
					Value: "\"\"\"\"\"\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind: EOF,
					Position: Position{
						Line:  1,
						Start: 6,
					},
				},
			},
		},
		{
			name: "simple string",
			src:  "\"\"\"simple string\"\"\"",
			want: []Token{
				{
					Kind:  BlockString,
					Value: "\"\"\"simple string\"\"\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind:  EOF,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 19,
					},
				},
			},
		},
		{
			name: "white space",
			src:  "\"\"\"  simple string  \"\"\"",
			want: []Token{
				{
					Kind:  BlockString,
					Value: "\"\"\"  simple string  \"\"\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind:  EOF,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 23,
					},
				},
			},
		},
		{
			name: "line feed",
			src:  "\"\"\" \nsimple string\"\"\"",
			want: []Token{
				{
					Kind:  BlockString,
					Value: "\"\"\" \nsimple string\"\"\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind:  EOF,
					Value: "",
					Position: Position{
						Line:  2,
						Start: 21,
					},
				},
			},
		},
		{
			name: "line carriage return",
			src:  "\"\"\" \rsimple string\"\"\"",
			want: []Token{
				{
					Kind:  BlockString,
					Value: "\"\"\" \rsimple string\"\"\"",
					Position: Position{
						Line:  1,
						Start: 1,
					},
				},
				{
					Kind:  EOF,
					Value: "",
					Position: Position{
						Line:  2,
						Start: 21,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(strings.NewReader(tt.src))

			gotTokens := make([]Token, 0)
			for {
				got, err := l.NextToken()
				if err != nil {
					t.Fatal(err)
				}

				gotTokens = append(gotTokens, got)
				if got.Kind == EOF || got.Kind == Invalid {
					break
				}
			}

			assert.Equal(t, tt.want, gotTokens)
		})
	}
}
