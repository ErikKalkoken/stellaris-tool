package parser

type Token uint

const (
	Illegal Token = iota
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
