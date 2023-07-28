package gogqllexer

type Kind int

const (
	Invalid Kind = iota
	EOF
	Name
	Bang
	Dollar
	Amp
	ParenL
	ParenR
	Spread
	Equal
	At
	Colon
	BracketL
	BracketR
	BraceL
	BraceR
	Pipe
	Comment
	Int
	Float
	String
)

type Position struct {
	Line  int
	Start int
}

type Token struct {
	Kind     Kind
	Value    string
	Position Position
}
