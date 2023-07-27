package gogqllexer

// https://spec.graphql.org/October2021/#sec-Int-Value
// https://spec.graphql.org/October2021/#sec-Float-Value
func isNumber(r rune) bool {
	switch r {
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	default:
		return false
	}
}

// https://spec.graphql.org/October2021/#ExponentPart
func isExponentPart(r rune) bool {
	switch r {
	case 'e', 'E':
		return true
	default:
		return false
	}
}

// https://spec.graphql.org/October2021/#FractionalPart
func isFractionalPart(r rune) bool {
	switch r {
	case '.':
		return true
	default:
		return false
	}
}

// https://spec.graphql.org/October2021/#Digit
func isDigit(r rune) bool {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	default:
		return false
	}
}

func isZero(r rune) bool {
	switch r {
	case '0':
		return true
	default:
		return false
	}
}

// https://spec.graphql.org/October2021/#NonZeroDigit
func isNonZeroDigit(r rune) bool {
	switch r {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	default:
		return false
	}
}

// https://spec.graphql.org/October2021/#NegativeSign
func isNegativeSign(r rune) bool {
	switch r {
	case '-':
		return true
	default:
		return false
	}
}
