package parser

type tokenType string

const (
	illegal       tokenType = "illegal"
	endOfFile     tokenType = "eof"
	equalSign     tokenType = "equalSign"
	bracketsOpen  tokenType = "bracketsOpen"
	bracketsClose tokenType = "bracketsClose"
	identifier    tokenType = "identifier"
	str           tokenType = "string"
	float         tokenType = "float"
	integer       tokenType = "integer"
	boolean       tokenType = "boolean"
)

type token struct {
	typ   tokenType
	value any
}
