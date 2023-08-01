package gogqllexer

import (
	"log"
	"strconv"
)

type Lexer struct {
	src            *Source
	line           int
	startByteIndex int
	endByteIndex   int
}

func New(src *Source) *Lexer {
	return &Lexer{
		src:  src,
		line: 1,
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
	// TODO: ignoreTokensまだまだある
	l.skipIgnoreTokens()
	l.startByteIndex = l.endByteIndex

	if l.endByteIndex >= len(l.src.Body) {
		return Token{
			Kind:  EOF,
			Value: "",
			Position: Position{
				Line:  l.line,
				Start: l.startByteIndex,
			},
		}, nil
	}

	// TODO: insignificant comma
	currentRune := rune(l.src.Body[l.startByteIndex])
	switch {
	case isNameStart(currentRune):
		kind, value, consumedByte, consumedLine := l.readNameToken()
		l.endByteIndex += consumedByte
		l.line += consumedLine

		return l.makeToken(kind, value), nil
	case isPunctuator(currentRune):
		kind, value, consumedByte, consumedLine := l.readPunctuatorToken()
		l.endByteIndex += consumedByte
		l.line += consumedLine

		return l.makeToken(kind, value), nil
	case isComment(currentRune):
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
	case isNumber(currentRune):
		kind, value, consumedByte, consumedLine := l.readNumber()
		l.endByteIndex += consumedByte
		l.line += consumedLine

		return l.makeToken(kind, value), nil
	case isStringValue(currentRune):
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

func (l *Lexer) readNumber() (kind Kind, value string, consumedByte int, consumedLine int) {
	isFloat := false
	if isNegativeSign(rune(l.src.Body[l.endByteIndex])) {
		consumedByte++
	}

	if isZero(rune(l.src.Body[l.endByteIndex+consumedByte])) {
		consumedByte++
		if l.endByteIndex+consumedByte < len(l.src.Body) && isDigit(rune(l.src.Body[l.endByteIndex+consumedByte])) {
			return Invalid, "", consumedByte, consumedLine
		}
	} else if isNonZeroDigit(rune(l.src.Body[l.endByteIndex+consumedByte])) {
		consumedByte++
		for l.endByteIndex+consumedByte < len(l.src.Body) {
			if isDigit(rune(l.src.Body[l.endByteIndex+consumedByte])) {
				consumedByte++
			} else {
				break
			}
		}
	} else {
		return Invalid, "", consumedByte, consumedLine
	}

	if l.endByteIndex+consumedByte < len(l.src.Body) && isFractionalPart(rune(l.src.Body[l.endByteIndex+consumedByte])) {
		consumedByte++
		isFloat = true
		for l.endByteIndex+consumedByte < len(l.src.Body) {
			if isDigit(rune(l.src.Body[l.endByteIndex+consumedByte])) {
				consumedByte++
			} else {
				break
			}
		}
	}

	if l.endByteIndex+consumedByte < len(l.src.Body) && isExponentPart(rune(l.src.Body[l.endByteIndex+consumedByte])) {
		consumedByte++
		isFloat = true
		for l.endByteIndex+consumedByte < len(l.src.Body) {
			if isDigit(rune(l.src.Body[l.endByteIndex+consumedByte])) || isSign(rune(l.src.Body[l.endByteIndex+consumedByte])) {
				consumedByte++
			} else {
				break
			}
		}
	}

	if isFloat {
		return Float, l.src.Body[l.startByteIndex : l.endByteIndex+consumedByte], consumedByte, consumedLine
	} else {
		return Int, l.src.Body[l.startByteIndex : l.endByteIndex+consumedByte], consumedByte, consumedLine
	}
}

func (l *Lexer) readPunctuatorToken() (kind Kind, value string, consumedByte int, consumedLine int) {
	consumedByte = 1
	switch l.src.Body[l.startByteIndex] {
	case '!':
		return Bang, "", consumedByte, consumedLine
	case '$':
		return Dollar, "", consumedByte, consumedLine
	case '&':
		return Amp, "", consumedByte, consumedLine
	case '(':
		return ParenL, "", consumedByte, consumedLine
	case ')':
		return ParenR, "", consumedByte, consumedLine
	case '.':
		if len(l.src.Body) <= l.startByteIndex+consumedByte+2 && l.src.Body[l.startByteIndex:l.startByteIndex+consumedByte+2] == "..." {
			consumedByte += 2
			return Spread, "", consumedByte, consumedLine
		}
		return Invalid, "", consumedByte, consumedLine
	case ':':
		return Colon, "", consumedByte, consumedLine
	case '=':
		return Equal, "", consumedByte, consumedLine
	case '@':
		return At, "", consumedByte, consumedLine
	case '[':
		return BracketL, "", consumedByte, consumedLine
	case ']':
		return BracketR, "", consumedByte, consumedLine
	case '{':
		return BraceL, "", consumedByte, consumedLine
	case '}':
		return BraceR, "", consumedByte, consumedLine
	case '|':
		return Pipe, "", consumedByte, consumedLine
	default:
		return Invalid, "", consumedByte, consumedLine
	}
}

func (l *Lexer) readNameToken() (kind Kind, value string, consumedByte int, consumedLine int) {
	for l.endByteIndex+consumedByte < len(l.src.Body) {
		if isNameContinue(rune(l.src.Body[l.endByteIndex+consumedByte])) {
			consumedByte++
		} else {
			break
		}
	}

	return Name, l.src.Body[l.startByteIndex : l.endByteIndex+consumedByte], consumedByte, consumedLine
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

// https://spec.graphql.org/October2021/#sec-Language.Source-Text.Ignored-Tokens
func (l *Lexer) skipIgnoreTokens() {
	for l.endByteIndex < len(l.src.Body) {
		r := rune(l.src.Body[l.endByteIndex])
		switch {
		case isWhiteSpace(r):
			l.endByteIndex++
		case isLineTerminator(r):
			l.line++
			l.endByteIndex++
			if l.endByteIndex < len(l.src.Body) && rune(l.src.Body[l.endByteIndex]) == '\n' {
				l.endByteIndex++
			}
		default:
			return
		}
	}
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
