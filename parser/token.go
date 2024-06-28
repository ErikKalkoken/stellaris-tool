package parser

type tokenType string

const (
	illegalType       tokenType = "illegal"
	eofType           tokenType = "eof"
	equalSignType     tokenType = "equalSign"
	bracketsOpenType  tokenType = "bracketsOpen"
	bracketsCloseType tokenType = "bracketsClose"
	identifierType    tokenType = "identifier"
	stringType        tokenType = "string"
	floatType         tokenType = "float"
	integerType       tokenType = "integer"
	booleanType       tokenType = "boolean"
	keywordType       tokenType = "keyword"
)

type token struct {
	typ   tokenType
	value any
}
