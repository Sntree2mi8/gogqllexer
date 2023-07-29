package gogqllexer

// https://spec.graphql.org/October2021/#sec-Comments
func isComment(r rune) bool {
	return r == '#'
}
