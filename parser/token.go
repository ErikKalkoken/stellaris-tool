package parser

type tokenType uint

const (
	illegalType tokenType = iota
	eofType
	whitespaceType
	equalSignType
	bracketsOpenType
	bracketsCloseType
	identifierType
	stringType
	floatType
	integerType
	booleanType
	keywordType
)

type token struct {
	typ   tokenType
	value any
}
