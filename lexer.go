package gogqllexer

import (
	"log"
)

type Lexer struct {
	// source
	src *Source

	// current line
	line int

	// ここから読む
	start int

	// ここまで読んだ
	end int
}

func New(src *Source) *Lexer {
	return &Lexer{
		src:  src,
		line: 1,
	}
}

func (l *Lexer) NextToken() (Token, error) {
	// TODO: ignoreTokensまだまだある
	l.ignoreTokens()
	l.start = l.end

	// 終端に達しているのでこれ以上Readするものがない
	if l.end >= len(l.src.Body) {
		return Token{
			Kind:  EOF,
			Value: "",
			Position: Position{
				Line:  l.line,
				Start: l.start,
			},
		}, nil
	}

	// TODO: insignificant comma
	// TODO: string
	currentRune := rune(l.src.Body[l.start])
	switch {
	case isNameStart(currentRune):
		for l.end < len(l.src.Body) {
			if isNameContinue(rune(l.src.Body[l.end])) {
				l.end++
			} else {
				break
			}
		}

		return Token{
			Kind:  Name,
			Value: l.src.Body[l.start:l.end],
			Position: Position{
				Line:  l.line,
				Start: l.start,
			},
		}, nil
	case isPunctuator(currentRune):
		switch l.src.Body[l.start] {
		case '!':
			l.end++
			return Token{
				Kind:  Bang,
				Value: "",
				Position: Position{
					Line:  l.line,
					Start: l.start,
				},
			}, nil
		case '$':
			l.end++
			return Token{
				Kind:  Dollar,
				Value: "",
				Position: Position{
					Line:  l.line,
					Start: l.start,
				},
			}, nil
		case '&':
			l.end++
			return Token{
				Kind:  Amp,
				Value: "",
				Position: Position{
					Line:  l.line,
					Start: l.start,
				},
			}, nil
		case '(':
			l.end++
			return Token{
				Kind:  ParenL,
				Value: "",
				Position: Position{
					Line:  l.line,
					Start: l.start,
				},
			}, nil
		case ')':
			l.end++
			return Token{
				Kind:  ParenR,
				Value: "",
				Position: Position{
					Line:  l.line,
					Start: l.start,
				},
			}, nil
		case '.':
			if len(l.src.Body) <= l.start+3 && l.src.Body[l.start:l.start+3] == "..." {
				l.end += 3
				return Token{
					Kind:  Spread,
					Value: "",
					Position: Position{
						Line:  l.line,
						Start: l.start,
					},
				}, nil
			}
		case ':':
			l.end++
			return Token{
				Kind:  Colon,
				Value: "",
				Position: Position{
					Line:  l.line,
					Start: l.start,
				},
			}, nil
		case '=':
			l.end++
			return Token{
				Kind:  Equal,
				Value: "",
				Position: Position{
					Line:  l.line,
					Start: l.start,
				},
			}, nil
		case '@':
			l.end++
			return Token{
				Kind:  At,
				Value: "",
				Position: Position{
					Line:  l.line,
					Start: l.start,
				},
			}, nil
		case '[':
			l.end++
			return Token{
				Kind:  BracketL,
				Value: "",
				Position: Position{
					Line:  l.line,
					Start: l.start,
				},
			}, nil
		case ']':
			l.end++
			return Token{
				Kind:  BracketR,
				Value: "",
				Position: Position{
					Line:  l.line,
					Start: l.start,
				},
			}, nil
		case '{':
			l.end++
			return Token{
				Kind:  BraceL,
				Value: "",
				Position: Position{
					Line:  l.line,
					Start: l.start,
				},
			}, nil
		case '}':
			l.end++
			return Token{
				Kind:  BraceR,
				Value: "",
				Position: Position{
					Line:  l.line,
					Start: l.start,
				},
			}, nil
		case '|':
			l.end++
			return Token{
				Kind:  Pipe,
				Value: "",
				Position: Position{
					Line:  l.line,
					Start: l.start,
				},
			}, nil
		}
	case isComment(currentRune):
		for l.end < len(l.src.Body) {
			if isLineTerminator(rune(l.src.Body[l.end])) {
				log.Println("line terminator")
				break
			} else {
				l.end++
			}
		}
		return Token{
			Kind:  Comment,
			Value: l.src.Body[l.start:l.end],
			Position: Position{
				Line:  l.line,
				Start: l.start,
			},
		}, nil
	case isNumber(currentRune):
		isFloat := false
		if isNegativeSign(rune(l.src.Body[l.end])) {
			l.end++
		}

		if isZero(rune(l.src.Body[l.end])) {
			l.end++
			if l.end < len(l.src.Body) && isDigit(rune(l.src.Body[l.end])) {
				return Token{
					Kind:  Invalid,
					Value: "invalid number token",
					Position: Position{
						Line:  l.line,
						Start: l.start,
					},
				}, nil
			}
		} else if isNonZeroDigit(rune(l.src.Body[l.end])) {
			l.end++
			for l.end < len(l.src.Body) {
				if isDigit(rune(l.src.Body[l.end])) {
					l.end++
				} else {
					break
				}
			}
		} else {
			return Token{
				Kind:  Invalid,
				Value: "invalid number token",
				Position: Position{
					Line:  l.line,
					Start: l.start,
				},
			}, nil
		}

		if isFractionalPart(rune(l.src.Body[l.end])) {
			l.end++
			isFloat = true
			for l.end < len(l.src.Body) {
				if isDigit(rune(l.src.Body[l.end])) {
					l.end++
				} else {
					break
				}
			}
		}

		if isExponentPart(rune(l.src.Body[l.end])) {
			l.end++
			isFloat = true
			for l.end < len(l.src.Body) {
				if isDigit(rune(l.src.Body[l.end])) {
					l.end++
				} else {
					break
				}
			}
		}

		if isFloat {
			return Token{
				Kind:  Float,
				Value: l.src.Body[l.start:l.end],
				Position: Position{
					Line:  l.line,
					Start: l.start,
				},
			}, nil
		} else {
			return Token{
				Kind:  Int,
				Value: l.src.Body[l.start:l.end],
				Position: Position{
					Line:  l.line,
					Start: l.start,
				},
			}, nil
		}
	}

	return Token{
		Kind:  Invalid,
		Value: "",
		Position: Position{
			Line:  l.line,
			Start: l.start,
		},
	}, nil
}

// https://spec.graphql.org/October2021/#NameStart
func isNameStart(r rune) bool {
	switch r {
	case '_', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
		return true
	default:
		return false
	}
}

// https://spec.graphql.org/October2021/#NameContinue
func isNameContinue(r rune) bool {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '_', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
		return true
	default:
		return false
	}
}

// https://spec.graphql.org/October2021/#sec-Punctuators
func isPunctuator(r rune) bool {
	switch r {
	case '!', '$', '&', '(', ')', '.', ':', '=', '@', '[', ']', '{', '}', '|':
		return true
	default:
		return false
	}
}

// https://spec.graphql.org/October2021/#sec-Comments
func isComment(r rune) bool {
	return r == '#'
}

// https://spec.graphql.org/October2021/#sec-Line-Terminators
func isLineTerminator(r rune) bool {
	switch r {
	case '\n', '\r':
		return true
	default:
		return false
	}
}

// https://spec.graphql.org/October2021/#sec-White-Space
func isWhiteSpace(r rune) bool {
	switch r {
	case ' ', '\t':
		return true
	default:
		return false
	}
}

// ignoreTokens ignore specific tokens
// https://spec.graphql.org/October2021/#sec-Language.Source-Text.Ignored-Tokens
func (l *Lexer) ignoreTokens() {
	for l.end < len(l.src.Body) {
		r := rune(l.src.Body[l.end])
		switch {
		case isWhiteSpace(r):
			l.end++
		case isLineTerminator(r):
			l.line++
			l.end++
			if l.end < len(l.src.Body) && rune(l.src.Body[l.end]) == '\n' {
				l.end++
			}
		default:
			return
		}
	}
}
