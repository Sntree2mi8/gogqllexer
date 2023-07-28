package gogqllexer

// https://spec.graphql.org/October2021/#sec-Punctuators
func isPunctuator(r rune) bool {
	switch r {
	case '!', '$', '&', '(', ')', '.', ':', '=', '@', '[', ']', '{', '}', '|':
		return true
	default:
		return false
	}
}
