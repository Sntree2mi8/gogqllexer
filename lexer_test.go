package gogqllexer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
						Start: 0,
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
						Start: 0,
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
						Start: 0,
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
						Start: 0,
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
						Start: 0,
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
						Start: 0,
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
						Start: 0,
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
						Start: 0,
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
						Start: 0,
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
						Start: 0,
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
						Start: 0,
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
						Start: 0,
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
						Start: 0,
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
						Start: 0,
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
					// TODO: 今は思いつくものがないのでエラーが起きたらfatalさせてしまう
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
					Kind:     Comment,
					Value:    "# This is comment",
					Position: Position{},
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

func TestLexer_NextToken(t *testing.T) {
	tests := []struct {
		name string
		src  *Source
		want []Token
	}{
		{
			name: "ignore whitespace",
			src: &Source{
				Body: "   \t",
				Name: "Spec_IgnoreWhiteSpace",
			},
			want: []Token{},
		},
		{
			name: "ignore whitespace",
			src: &Source{
				Body: "   type",
				Name: "Spec_IgnoreWhiteSpace",
			},
			want: []Token{
				{
					Kind:  Name,
					Value: "type",
					Position: Position{
						Line:  1,
						Start: 3,
					},
				},
			},
		},
		{
			name: "punctuator bang",
			src: &Source{
				Body: "!",
				Name: "SpecPunctuatorBang",
			},
			want: []Token{
				{
					Kind:  Bang,
					Value: "",
					Position: Position{
						Line:  1,
						Start: 0,
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
					// TODO: 今は思いつくものがないのでエラーが起きたらfatalさせてしまう
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
