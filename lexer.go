package gogqllexer

import (
	"io"
	"log"
)

type Lexer struct {
	io.RuneScanner

	src            *Source
	line           int
	startByteIndex int
	endByteIndex   int
}

func New(src *Source, scanner io.RuneScanner) *Lexer {
	return &Lexer{
		RuneScanner:    scanner,
		src:            src,
		line:           1,
		startByteIndex: 0,
		endByteIndex:   0,
	}
}

func (l *Lexer) makeEOFToken() Token {
	return Token{
		Kind:  EOF,
		Value: "",
		Position: Position{
			Line:  l.line,
			Start: l.startByteIndex,
		},
	}
}

func (l *Lexer) makeToken(kind Kind, value string) Token {
	return Token{
		Kind:  kind,
		Value: value,
		Position: Position{
			Line:  l.line,
			Start: l.startByteIndex + 1,
		},
	}
}

func (l *Lexer) NextToken() (Token, error) {
	// skip ignore tokens
	for {
		r, s, err := l.ReadRune()
		if err != nil {
			return l.makeEOFToken(), nil
		}

		// TODO: more ignore tokens
		switch {
		case isWhiteSpace(r):
			l.startByteIndex += s
			continue
		case isLineTerminator(r):
			l.startByteIndex += s
			l.line++
			r, s, err = l.ReadRune()
			if err != nil {
				return l.makeEOFToken(), nil
			}
			if r == '\n' {
				l.startByteIndex += s
			} else {
				if err = l.UnreadRune(); err != nil {
					return l.makeToken(Invalid, ""), err
				}
			}
			continue
		default:
			if err = l.UnreadRune(); err != nil {
				return l.makeToken(Invalid, ""), err
			}
		}
		break
	}

	// TODO: insignificant comma
	r, err := l.peek()
	if err != nil {
		return l.makeToken(Invalid, ""), err
	}
	switch {
	case isNameStart(r):
		t, consumedByte := l.readNameToken()
		l.startByteIndex += consumedByte
		return t, nil
	case isPunctuator(r):
		t, consumedByte := l.readPunctuatorToken()
		l.startByteIndex += consumedByte
		return t, nil
	case isNumber(r):
		t, consumedByte := l.readNumber()
		l.startByteIndex += consumedByte
		return t, nil
	case isStringValue(r):
		t, consumedByte, consumedLine := l.readStringToken()
		l.startByteIndex += consumedByte
		l.line += consumedLine
		return t, nil
	case isComment(r):
		for l.endByteIndex < len(l.src.Body) {
			if isLineTerminator(rune(l.src.Body[l.endByteIndex])) {
				log.Println("line terminator")
				break
			} else {
				l.endByteIndex++
			}
		}
		return Token{
			Kind:  Comment,
			Value: l.src.Body[l.startByteIndex:l.endByteIndex],
			Position: Position{
				Line:  l.line,
				Start: l.startByteIndex,
			},
		}, nil
	}

	return l.makeToken(Invalid, ""), nil
}

func (l *Lexer) peek() (rune, error) {
	r, _, err := l.ReadRune()
	if err != nil {
		return 0, err
	}
	_ = l.UnreadRune()

	return r, nil
}

func (l *Lexer) readNumber() (token Token, consumedByte int) {
	isFloat := false

	r, s, err := l.ReadRune()
	if err != nil {
		return l.makeEOFToken(), consumedByte
	}
	consumedByte += s

	if isNegativeSign(r) {
		r, s, err = l.ReadRune()
		if err != nil {
			return l.makeToken(Invalid, ""), consumedByte
		}
		consumedByte += s
	}

	if isZero(r) {
		r, s, err = l.ReadRune()
		if err != nil {
			return l.makeToken(Int, l.src.Body[l.startByteIndex:l.startByteIndex+consumedByte]), consumedByte
		}
		consumedByte += s
		if isDigit(r) || isNameStart(r) && !isExponentPart(r) {
			return l.makeToken(Invalid, ""), consumedByte
		}
	} else if isNonZeroDigit(r) {
		for {
			r, s, err = l.ReadRune()
			if err != nil {
				return l.makeToken(Int, l.src.Body[l.startByteIndex:l.startByteIndex+consumedByte]), consumedByte
			}
			consumedByte += s
			if isDigit(r) {
				continue
			} else if isNameStart(r) && !isExponentPart(r) {
				return l.makeToken(Invalid, ""), consumedByte
			} else {
				break
			}
		}
	} else {
		return l.makeToken(Invalid, ""), consumedByte
	}

	if isFractionalPart(r) {
		isFloat = true
		// fractional part must be followed by at least one digit
		r, s, err = l.ReadRune()
		if err != nil {
			return l.makeToken(Invalid, ""), consumedByte
		}
		consumedByte += s
		if !isDigit(r) {
			return l.makeToken(Invalid, ""), consumedByte
		}

		for {
			r, s, err = l.ReadRune()
			if err != nil {
				break
			}
			consumedByte += s
			if isDigit(r) {
				continue
			} else if (isNameStart(r) && !isExponentPart(r)) || r == '.' {
				return l.makeToken(Invalid, ""), consumedByte
			} else {
				break
			}
		}
	}

	if isExponentPart(r) {
		isFloat = true

		// check opt sign
		r, err = l.peek()
		if err != nil {
			return l.makeToken(Invalid, ""), consumedByte
		}
		if isSign(r) {
			r, s, err = l.ReadRune()
			if err != nil {
				return l.makeToken(Invalid, ""), consumedByte
			}
			consumedByte += s
		}

		for {
			r, s, err = l.ReadRune()
			if err != nil {
				break
			}
			consumedByte += s
			if isDigit(r) {
				continue
			} else if isNameStart(r) || r == '.' {
				return l.makeToken(Invalid, ""), consumedByte
			} else {
				_ = l.UnreadRune()
				consumedByte -= s
				break
			}
		}
	}

	if isFloat {
		return l.makeToken(Float, l.src.Body[l.startByteIndex:l.startByteIndex+consumedByte]), consumedByte
	} else {
		return l.makeToken(Int, l.src.Body[l.startByteIndex:l.startByteIndex+consumedByte]), consumedByte
	}
}

func (l *Lexer) readPunctuatorToken() (token Token, consumedByte int) {
	r, consumedByte, err := l.ReadRune()
	if err != nil {
		return l.makeEOFToken(), consumedByte
	}

	switch r {
	case '!':
		return l.makeToken(Bang, ""), consumedByte
	case '$':
		return l.makeToken(Dollar, ""), consumedByte
	case '&':
		return l.makeToken(Amp, ""), consumedByte
	case '(':
		return l.makeToken(ParenL, ""), consumedByte
	case ')':
		return l.makeToken(ParenR, ""), consumedByte
	case '.':
		for i := 0; i < 2; i++ {
			r, s, err := l.ReadRune()
			if err != nil {
				return l.makeToken(Invalid, ""), consumedByte
			}
			consumedByte += s
			if r != '.' {
				return l.makeToken(Invalid, ""), consumedByte
			}
		}
		return l.makeToken(Spread, ""), consumedByte
	case ':':
		return l.makeToken(Colon, ""), consumedByte
	case '=':
		return l.makeToken(Equal, ""), consumedByte
	case '@':
		return l.makeToken(At, ""), consumedByte
	case '[':
		return l.makeToken(BracketL, ""), consumedByte
	case ']':
		return l.makeToken(BracketR, ""), consumedByte
	case '{':
		return l.makeToken(BraceL, ""), consumedByte
	case '}':
		return l.makeToken(BraceR, ""), consumedByte
	case '|':
		return l.makeToken(Pipe, ""), consumedByte
	default:
		return l.makeToken(Invalid, ""), consumedByte
	}
}

func (l *Lexer) readNameToken() (token Token, consumedByte int) {
	for {
		r, s, err := l.ReadRune()
		if err != nil {
			//EOF
			return l.makeToken(Name, l.src.Body[l.startByteIndex:l.startByteIndex+consumedByte]), consumedByte
		}
		if isNameContinue(r) {
			consumedByte += s
			continue
		}
		if err = l.UnreadRune(); err != nil {
			return l.makeToken(Invalid, ""), consumedByte
		}
		return l.makeToken(Name, l.src.Body[l.startByteIndex:l.startByteIndex+consumedByte]), consumedByte
	}
}

func (l *Lexer) readStringToken() (token Token, consumedByte int, consumedLine int) {
	r, s, err := l.ReadRune()
	if err != nil {
		return l.makeEOFToken(), consumedByte, consumedLine
	}
	consumedByte += s

	if r != '"' {
		return l.makeToken(Invalid, ""), consumedByte, consumedLine
	}

	isBlockString := false

	makeStringToken := func() Token {
		return Token{
			Kind:  String,
			Value: l.src.Body[l.startByteIndex : l.startByteIndex+consumedByte],
			Position: Position{
				Line:  l.line + consumedLine,
				Start: l.startByteIndex + 1,
			},
		}
	}

	makeBlockStringToken := func() Token {
		return Token{
			Kind:  BlockString,
			Value: l.src.Body[l.startByteIndex : l.startByteIndex+consumedByte],
			Position: Position{
				Line:  l.line + consumedLine,
				Start: l.startByteIndex + 1,
			},
		}
	}

StringReadLoop:
	for {
		r, s, err = l.ReadRune()
		if err != nil {
			return l.makeToken(Invalid, ""), consumedByte, consumedLine
		}
		consumedByte += s

		switch r {
		case '\n', '\r':
			return l.makeToken(Invalid, ""), consumedByte, consumedLine
		case '"':
			r, err = l.peek()
			if err != nil {
				return makeStringToken(), consumedByte, consumedLine
			}
			if r == '"' {
				isBlockString = true
				_, s, _ = l.ReadRune()
				consumedByte += s
				break StringReadLoop
			} else {
				return makeStringToken(), consumedByte, consumedLine
			}
		case '\\':
			r, s, err = l.ReadRune()
			if err != nil {
				return l.makeToken(Invalid, ""), consumedByte, consumedLine
			}
			consumedByte += s

			switch r {
			default:
				return l.makeToken(Invalid, ""), consumedByte, consumedLine
			case 'u':
				for i := 0; i < 4; i++ {
					r, s, err = l.ReadRune()
					if err != nil {
						return l.makeToken(Invalid, ""), consumedByte, consumedLine
					}
					consumedByte += s

					if !isHexDigit(r) {
						return l.makeToken(Invalid, ""), consumedByte, consumedLine
					}
				}
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
				break
			}
		default:
			if r < 0x0020 && r != '\t' {
				return l.makeToken(Invalid, ""), consumedByte, consumedLine
			}
		}
	}

	if isBlockString {
		for {
			r, s, err = l.ReadRune()
			if err != nil {
				return l.makeToken(Invalid, ""), consumedByte, consumedLine
			}
			consumedByte += s

			switch r {
			case '\n':
				consumedLine++
			case '\r':
				consumedLine++
				if r, err = l.peek(); err != nil {
					return l.makeToken(Invalid, ""), consumedByte, consumedLine
				} else if r == '\n' {
					_, s, _ = l.ReadRune()
					consumedByte += s
				}
			case '"':
				for i := 0; i < 2; i++ {
					r, s, err = l.ReadRune()
					if err != nil {
						return l.makeToken(Invalid, ""), consumedByte, consumedLine
					}
					consumedByte += s
					if r != '"' {
						return l.makeToken(Invalid, ""), consumedByte, consumedLine
					}
				}
				return makeBlockStringToken(), consumedByte, consumedLine
			case '\\':
				r, s, err = l.ReadRune()
				if err != nil {
					return l.makeToken(Invalid, ""), consumedByte, consumedLine
				}
				consumedByte += s

				switch r {
				default:
					return l.makeToken(Invalid, ""), consumedByte, consumedLine
				case 'u':
					for i := 0; i < 4; i++ {
						r, s, err = l.ReadRune()
						if err != nil {
							return l.makeToken(Invalid, ""), consumedByte, consumedLine
						}
						consumedByte += s

						if !isHexDigit(r) {
							return l.makeToken(Invalid, ""), consumedByte, consumedLine
						}
					}
				case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
					break
				}
			default:
				if r < 0x0020 && r != '\t' && r != '\n' && r != '\r' {
					return l.makeToken(Invalid, ""), consumedByte, consumedLine
				}
			}
		}
	}

	return l.makeToken(Invalid, ""), consumedByte, consumedLine
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
