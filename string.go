package gogqllexer

// https://spec.graphql.org/draft/#sec-String-Value
func isStringValue(r rune) bool {
	return r == '"'
}
