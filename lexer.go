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
		case isComma(r):
			consumedByte += s
			continue
		case isUnicodeBOM(r):
			consumedByte += s
			continue
		case isWhiteSpace(r):
			consumedByte += s
			continue
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
		default:
			_ = l.UnreadRune()
			break ReadIgnoredTokenLoop
		}
	}

	return consumedByte, consumedLine
}

func (l *Lexer) NextToken() (Token, error) {
	consumedByte, consumedLine := l.skipIgnoreTokens()
	l.startByteIndex += consumedByte
	l.line += consumedLine

	r, err := l.peek()
	if err != nil {
		return l.makeEOFToken(), nil
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
		t, consumedByte := l.readComment()
		l.startByteIndex += consumedByte
		return t, nil
	default:
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

func (l *Lexer) readComment() (token Token, consumedByte int) {
	value := make([]rune, 0)

	r, s, err := l.ReadRune()
	if err != nil {
		return l.makeEOFToken(), consumedByte
	}
	consumedByte += s
	value = append(value, r)

	if r != '#' {
		return l.makeToken(Invalid, ""), consumedByte
	}

ReadCommentLoop:
	for {
		r, err = l.peek()
		if err != nil {
			break
		}

		switch {
		case isLineTerminator(r), r < 0x0020 && r != '\t':
			break ReadCommentLoop
		default:
			r, s, _ = l.ReadRune()
			consumedByte += s
			value = append(value, r)
		}
	}

	return Token{
		Kind:  Comment,
		Value: string(value),
		Position: Position{
			Line:  l.line,
			Start: l.startByteIndex,
		},
	}, consumedByte
}

func (l *Lexer) readNumber() (token Token, consumedByte int) {
	isFloat := false
	value := make([]rune, 0)

	r, s, err := l.ReadRune()
	if err != nil {
		return l.makeEOFToken(), consumedByte
	}
	consumedByte += s
	value = append(value, r)

	if isNegativeSign(r) {
		r, s, err = l.ReadRune()
		if err != nil {
			return l.makeToken(Invalid, ""), consumedByte
		}
		consumedByte += s
		value = append(value, r)
	}

	if isZero(r) {
		r, s, err = l.ReadRune()
		if err != nil {
			return l.makeToken(Int, string(value)), consumedByte
		}
		consumedByte += s
		value = append(value, r)

		if isDigit(r) || isNameStart(r) && !isExponentPart(r) {
			return l.makeToken(Invalid, ""), consumedByte
		}
	} else if isNonZeroDigit(r) {
		for {
			r, s, err = l.ReadRune()
			if err != nil {
				return l.makeToken(Int, string(value)), consumedByte
			}
			consumedByte += s
			value = append(value, r)

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
		value = append(value, r)

		if !isDigit(r) {
			return l.makeToken(Invalid, ""), consumedByte
		}

		for {
			r, s, err = l.ReadRune()
			if err != nil {
				break
			}
			consumedByte += s
			value = append(value, r)

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
			value = append(value, r)
		}

		for {
			r, s, err = l.ReadRune()
			if err != nil {
				break
			}
			consumedByte += s
			value = append(value, r)

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
		return l.makeToken(Float, string(value)), consumedByte
	} else {
		return l.makeToken(Int, string(value)), consumedByte
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

// https://spec.graphql.org/October2021/#sec-White-Space
func isWhiteSpace(r rune) bool {
	switch r {
	case ' ', '\t':
		return true
	default:
		return false
	}
}

// https://spec.graphql.org/October2021/#sec-Unicode
func isUnicodeBOM(r rune) bool {
	switch r {
	case '\uFEFF':
		return true
	default:
		return false
	}
}

// https://spec.graphql.org/October2021/#sec-Insignificant-Commas
func isComma(r rune) bool {
	return r == ','
}
