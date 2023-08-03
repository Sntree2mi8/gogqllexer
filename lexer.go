package gogqllexer

import (
	"io"
	"log"
	"strconv"
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
	case isStringValue(r):
		if l.endByteIndex+3 < len(l.src.Body) && l.src.Body[l.endByteIndex:l.endByteIndex+3] == `"""` {
			kind, value, consumedByte, consumedLine := l.readStringBlockToken()
			l.endByteIndex += consumedByte
			l.line += consumedLine
			return l.makeToken(kind, value), nil
		} else {
			kind, value, consumedByte, consumedLine := l.readStringToken()
			l.endByteIndex += consumedByte
			l.line += consumedLine
			return l.makeToken(kind, value), nil
		}
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
			} else {
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

func (l *Lexer) readStringToken() (kind Kind, value string, consumedByte int, consumedLine int) {
	// consumedByte initial value is 1 because of skipping double quote
	consumedByte = 1

StringReadLoop:
	for l.endByteIndex+consumedByte < len(l.src.Body) {
		switch rune(l.src.Body[l.endByteIndex+consumedByte]) {
		case '\n', '\r':
			consumedByte++
			break StringReadLoop
		case '"':
			consumedByte++
			return String, l.src.Body[l.startByteIndex : l.endByteIndex+consumedByte], consumedByte, consumedLine
		case '\\':
			consumedByte++
			if l.endByteIndex+consumedByte < len(l.src.Body) {
				nextRune := rune(l.src.Body[l.endByteIndex+consumedByte])
				switch nextRune {
				case 'u':
					consumedByte++
					if l.endByteIndex+consumedByte+4 >= len(l.src.Body) {
						break StringReadLoop
					}
					_, err := strconv.ParseUint(l.src.Body[l.endByteIndex+consumedByte:l.endByteIndex+consumedByte+4], 16, 64)
					if err != nil {
						break StringReadLoop
					}
					consumedByte += 4
				case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
					consumedByte++
				default:
					consumedByte++
					break StringReadLoop
				}
			} else {
				break StringReadLoop
			}
		default:
			if rune(l.src.Body[l.endByteIndex+consumedByte]) < 0x0020 && rune(l.src.Body[l.endByteIndex+consumedByte]) != '\t' {
				consumedByte++
				break StringReadLoop
			}
			consumedByte++
		}
	}

	return Invalid, "", consumedByte, consumedLine
}

func (l *Lexer) readStringBlockToken() (kind Kind, value string, consumedByte int, consumedLine int) {
	// consumedByte initial value is 3 because of skipping triple double quote
	consumedByte = 3

BlockStringReadLoop:
	for l.endByteIndex+consumedByte < len(l.src.Body) {
		switch rune(l.src.Body[l.endByteIndex+consumedByte]) {
		case '\n':
			consumedByte++
			consumedLine++
		case '\r':
			consumedByte++
			consumedLine++
			if l.endByteIndex+consumedByte < len(l.src.Body) && rune(l.src.Body[l.endByteIndex+consumedByte]) == '\n' {
				consumedByte++
			}
		case '"':
			if l.endByteIndex+consumedByte+3 <= len(l.src.Body) && l.src.Body[l.endByteIndex+consumedByte:l.endByteIndex+consumedByte+3] == `"""` {
				consumedByte += 3
				return BlockString, l.src.Body[l.startByteIndex : l.endByteIndex+consumedByte], consumedByte, consumedLine
			} else {
				consumedByte++
			}
		case '\\':
			consumedByte++
			if l.endByteIndex+consumedByte < len(l.src.Body) {
				nextRune := rune(l.src.Body[l.endByteIndex+consumedByte])
				switch nextRune {
				case 'u':
					consumedByte++
					if l.endByteIndex+consumedByte+4 >= len(l.src.Body) {
						break BlockStringReadLoop
					}
					_, err := strconv.ParseUint(l.src.Body[l.endByteIndex+consumedByte:l.endByteIndex+consumedByte+4], 16, 64)
					if err != nil {
						break BlockStringReadLoop
					}
					consumedByte += 4
				case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
					consumedByte++
				default:
					consumedByte++
					break BlockStringReadLoop
				}
			} else {
				break BlockStringReadLoop
			}
		default:
			r := rune(l.src.Body[l.endByteIndex+consumedByte])
			if r < 0x0020 && r != '\t' && r != '\n' && r != '\r' {
				consumedByte++
				break BlockStringReadLoop
			}
			consumedByte++
		}
	}

	return Invalid, "", consumedByte, consumedLine
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
