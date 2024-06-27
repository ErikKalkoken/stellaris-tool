package parser

type tokenType uint

const (
	Illegal tokenType = iota
	Eof
	Whitespace
	Equal
	BracketsOpen
	BracketsClose
	Identifier
	String
	Float
	Integer
	Boolean
)

type Token struct {
	typ   tokenType
	value any
}
