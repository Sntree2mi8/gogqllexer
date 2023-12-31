package gogqllexer

import (
	"io"
)

type Lexer struct {
	io.RuneScanner

	line           int
	startByteIndex int
}

func New(scanner io.RuneScanner) *Lexer {
	return &Lexer{
		RuneScanner:    scanner,
		line:           1,
		startByteIndex: 0,
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

func (l *Lexer) NextToken() Token {
	consumedByte, consumedLine := l.skipIgnoreTokens()
	l.startByteIndex += consumedByte
	l.line += consumedLine

	r, err := l.peek()
	if err != nil {
		return l.makeEOFToken()
	}
	switch {
	case isNameStart(r):
		t, consumedByte := l.readNameToken()
		l.startByteIndex += consumedByte
		return t
	case isPunctuator(r):
		t, consumedByte := l.readPunctuatorToken()
		l.startByteIndex += consumedByte
		return t
	case isNumber(r):
		t, consumedByte := l.readNumber()
		l.startByteIndex += consumedByte
		return t
	case isStringValue(r):
		t, consumedByte, consumedLine := l.readStringToken()
		l.startByteIndex += consumedByte
		l.line += consumedLine
		return t
	default:
	}

	return l.makeToken(Invalid, "")
}

func (l *Lexer) peek() (rune, error) {
	r, _, err := l.ReadRune()
	if err != nil {
		return 0, err
	}
	_ = l.UnreadRune()

	return r, nil
}

func isNumber(r rune) bool {
	return r == '-' || isDigit(r)
}

// https://spec.graphql.org/October2021/#ExponentPart
func isExponentPart(r rune) bool {
	return r == 'e' || r == 'E'
}

// https://spec.graphql.org/October2021/#FractionalPart
func isFractionalPart(r rune) bool {
	return r == '.'
}

// https://spec.graphql.org/October2021/#Digit
func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isHexDigit(r rune) bool {
	return isDigit(r) || ('a' <= r && r <= 'f') || ('A' <= r && r <= 'F')
}

func (l *Lexer) readNumber() (token Token, consumedByte int) {
	isFloat := false
	// TODO: fix name. スコープがでかいのに不鮮明な名前
	s := 0
	value := make([]rune, 0)

	// check EOF
	r, err := l.peek()
	if err != nil {
		return l.makeEOFToken(), consumedByte
	}

	// check sign
	if r == '-' {
		r, s, _ = l.ReadRune()
		consumedByte += s
		value = append(value, r)

		r, err = l.peek()
		if err != nil {
			return l.makeToken(Invalid, ""), consumedByte
		}
	}

	// read integer part
	leadingZero := false

	// check reading 0
	if r == '0' {
		leadingZero = true
		r, s, _ = l.ReadRune()
		consumedByte += s
		value = append(value, r)
	}

	for {
		r, err = l.peek()
		if err != nil {
			return l.makeToken(Int, string(value)), consumedByte
		}

		if isDigit(r) {
			if leadingZero {
				return l.makeToken(Invalid, ""), consumedByte
			}
			r, s, _ = l.ReadRune()
			consumedByte += s
			value = append(value, r)
		} else if isNameStart(r) && !isExponentPart(r) {
			_, s, _ = l.ReadRune()
			consumedByte += s
			return l.makeToken(Invalid, ""), consumedByte
		} else {
			break
		}
	}

	if isFractionalPart(r) {
		isFloat = true
		r, s, _ = l.ReadRune()
		consumedByte += s
		value = append(value, r)

		// dot must be followed by at least one digit
		r, s, err = l.ReadRune()
		if err != nil {
			return l.makeToken(Invalid, ""), consumedByte
		}
		if !isDigit(r) {
			return l.makeToken(Invalid, ""), consumedByte
		}
		consumedByte += s
		value = append(value, r)

		for {
			r, err = l.peek()
			if err != nil {
				return l.makeToken(Float, string(value)), consumedByte
			}

			if isDigit(r) {
				r, s, _ = l.ReadRune()
				consumedByte += s
				value = append(value, r)
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
		r, s, _ = l.ReadRune()
		consumedByte += s
		value = append(value, r)

		// check opt sign
		r, err = l.peek()
		if err != nil {
			return l.makeToken(Invalid, ""), consumedByte
		}
		if r == '-' || r == '+' {
			r, s, _ = l.ReadRune()
			consumedByte += s
			value = append(value, r)
		}

		// must be followed by at least one digit
		r, s, err = l.ReadRune()
		if err != nil {
			return l.makeToken(Invalid, ""), consumedByte
		}
		if !isDigit(r) {
			return l.makeToken(Invalid, ""), consumedByte
		}
		consumedByte += s
		value = append(value, r)

		for {
			r, err = l.peek()
			if err != nil {
				return l.makeToken(Float, string(value)), consumedByte
			}

			if isDigit(r) {
				r, s, _ = l.ReadRune()
				consumedByte += s
				value = append(value, r)

				continue
			} else if isNameStart(r) || r == '.' {
				return l.makeToken(Invalid, ""), consumedByte
			} else {
				break
			}
		}
	}

	if isFloat {
		return l.makeToken(Float, string(value)), consumedByte
	} else {
		return l.makeToken(Int, string(value)), consumedByte
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

// https://spec.graphql.org/October2021/#NameStart
func isNameStart(r rune) bool {
	switch {
	case r == '_' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z':
		return true
	default:
		return false
	}
}

// https://spec.graphql.org/October2021/#NameContinue
func isNameContinue(r rune) bool {
	return isNameStart(r) || isDigit(r)
}

func (l *Lexer) readNameToken() (token Token, consumedByte int) {
	value := make([]rune, 0)
	for {
		r, s, err := l.ReadRune()
		if err != nil {
			//EOF
			return l.makeToken(Name, string(value)), consumedByte
		}
		if isNameContinue(r) {
			consumedByte += s
			value = append(value, r)
			continue
		}
		_ = l.UnreadRune()

		return l.makeToken(Name, string(value)), consumedByte
	}
}

// https://spec.graphql.org/draft/#sec-String-Value
func isStringValue(r rune) bool {
	return r == '"'
}

func (l *Lexer) readStringToken() (token Token, consumedByte int, consumedLine int) {
	value := make([]rune, 0)
	r, s, err := l.ReadRune()
	if err != nil {
		return l.makeEOFToken(), consumedByte, consumedLine
	}
	consumedByte += s
	value = append(value, r)

	if r != '"' {
		return l.makeToken(Invalid, ""), consumedByte, consumedLine
	}

	isBlockString := false

	makeStringToken := func(v string) Token {
		return Token{
			Kind:  String,
			Value: v,
			Position: Position{
				Line:  l.line + consumedLine,
				Start: l.startByteIndex + 1,
			},
		}
	}

	makeBlockStringToken := func(v string) Token {
		return Token{
			Kind:  BlockString,
			Value: v,
			Position: Position{
				Line:  l.line,
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
		value = append(value, r)

		switch r {
		case '\n', '\r':
			return l.makeToken(Invalid, ""), consumedByte, consumedLine
		case '"':
			r, err = l.peek()
			if err != nil {
				return makeStringToken(string(value)), consumedByte, consumedLine
			}
			if r == '"' {
				isBlockString = true
				r, s, _ = l.ReadRune()
				consumedByte += s
				value = append(value, r)
				break StringReadLoop
			} else {
				return makeStringToken(string(value)), consumedByte, consumedLine
			}
		case '\\':
			r, s, err = l.ReadRune()
			if err != nil {
				return l.makeToken(Invalid, ""), consumedByte, consumedLine
			}
			consumedByte += s
			value = append(value, r)

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
					value = append(value, r)

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
			value = append(value, r)

			switch r {
			case '\n':
				consumedLine++
			case '\r':
				consumedLine++
				if r, err = l.peek(); err != nil {
					return l.makeToken(Invalid, ""), consumedByte, consumedLine
				} else if r == '\n' {
					r, s, _ = l.ReadRune()
					consumedByte += s
					value = append(value, r)
				}
			case '"':
				for i := 0; i < 2; i++ {
					r, s, err = l.ReadRune()
					if err != nil {
						return l.makeToken(Invalid, ""), consumedByte, consumedLine
					}
					consumedByte += s
					value = append(value, r)
					if r != '"' {
						return l.makeToken(Invalid, ""), consumedByte, consumedLine
					}
				}
				return makeBlockStringToken(string(value)), consumedByte, consumedLine
			case '\\':
				r, s, err = l.ReadRune()
				if err != nil {
					return l.makeToken(Invalid, ""), consumedByte, consumedLine
				}
				consumedByte += s
				value = append(value, r)

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
						value = append(value, r)

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

// https://spec.graphql.org/October2021/#sec-Language.Source-Text.Ignored-Tokens
// Commentを除く
func (l *Lexer) skipIgnoreTokens() (consumedByte int, consumedLine int) {
ReadIgnoredTokenLoop:
	for {
		r, s, err := l.ReadRune()
		if err != nil {
			break ReadIgnoredTokenLoop
		}

		switch {
		// insignificant comma
		case r == ',':
			consumedByte += s
			continue
		// unicode BOM
		case r == '\uFEFF':
			consumedByte += s
			continue
		// whitespace
		case r == ' ', r == '\t':
			consumedByte += s
			continue
		// line terminator
		case isLineTerminator(r):
			consumedByte += s
			consumedLine++
			r, s, err = l.ReadRune()
			if err != nil {
				break ReadIgnoredTokenLoop
			}
			if r == '\n' {
				consumedByte += s
			} else {
				_ = l.UnreadRune()
			}
			continue
		case r == '#':
			consumedByte += s
			for {
				r, err = l.peek()
				if err != nil {
					break
				}

				if isLineTerminator(r) || r < 0x0020 && r != '\t' {
					break
				}

				_, s, _ = l.ReadRune()
				consumedByte += s
			}
		default:
			_ = l.UnreadRune()
			break ReadIgnoredTokenLoop
		}
	}

	return consumedByte, consumedLine
}
