package parser

import (
	"fmt"
	"io"
	"strconv"
)

type Keyword uint

// Keywords
const (
	None Keyword = iota + 1
	NotSet
	Indeterminable
	Male
	Female
)

// Parser represents a parser.
type Parser struct {
	// Provides a stream of tokens
	lex *Lexer
	// Stack of latest tokens so we can go back
	ts stack[token]
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{lex: NewLexer(r), ts: newStack[token](3)}
}
func (p *Parser) Parse() (map[string]any, error) {
	x := make(map[string]any)

loop:
	for {
		var key string
		var value any

		// First token should be identifier or integer
		switch tok := p.nextRegularToken(); tok.typ {
		case eofType, bracketsCloseType:
			break loop
		case identifierType:
			key = tok.value.(string)
		case integerType:
			key = strconv.Itoa(tok.value.(int))
		default:
			return nil, p.makeError("found %v, expected identifier or integer", tok)
		}

		// Next should be an equal sign
		if tok := p.nextRegularToken(); tok.typ != equalSignType {
			return nil, p.makeError("found %v, expected equal sign", tok)
		}

		// Next should be some kind of value
		switch tok := p.nextRegularToken(); tok.typ {
		case stringType, floatType, integerType, booleanType, keywordType:
			value = tok.value
		case bracketsOpenType:
			tok2 := p.nextRegularToken()
			switch tok2.typ {
			case bracketsCloseType:
				// Empty object
				value = struct{}{}
			case bracketsOpenType:
				// Array of objects
				x := make([]map[string]any, 0)
				for {
					v2, err := p.Parse()
					if err != nil {
						return nil, err
					}
					x = append(x, v2)
					tok3 := p.nextRegularToken()
					if tok3.typ == bracketsCloseType {
						break
					} else if tok3.typ != bracketsOpenType {
						return nil, p.makeError("Unexpected token %v in obj array", tok3)
					}
				}
				value = x
			case identifierType:
				// Normal object
				p.backup(tok2)
				x, err := p.Parse()
				if err != nil {
					return nil, err
				}
				value = x
			case integerType:
				tok3 := p.nextRegularToken()
				p.backup(tok3)
				p.backup(tok2)
				if tok3.typ == equalSignType {
					// ID object
					x, err := p.Parse()
					if err != nil {
						return nil, err
					}
					value = x
					break
				}
				// Array of integer
				x := make([]int, 0)
				for {
					tok3 := p.nextRegularToken()
					if tok3.typ == bracketsCloseType {
						value = x
						break
					}
					y, ok := tok3.value.(int)
					if !ok {
						return nil, p.makeError("Expected type integer, but got: %v", tok3)
					}
					x = append(x, y)
				}
			case floatType:
				// Array of float
				p.backup(tok2)
				x := make([]float64, 0)
				for {
					tok3 := p.nextRegularToken()
					if tok3.typ == bracketsCloseType {
						value = x
						break
					}
					y, ok := tok3.value.(float64)
					if !ok {
						return nil, p.makeError("Expected type float, but got: %v", tok2)
					}
					x = append(x, y)
				}
			case stringType:
				// Array of string
				p.backup(tok2)
				x := make([]string, 0)
				for {
					tok3 := p.nextRegularToken()
					if tok3.typ == bracketsCloseType {
						value = x
						break
					}
					y, ok := tok3.value.(string)
					if !ok {
						return nil, p.makeError("Expected type string, but got: %v", tok3)
					}
					x = append(x, y)
				}
			default:
				return nil, p.makeError("invalid token %v for array", tok2)
			}

		default:
			return nil, p.makeError("found %v, expected a value", tok)
		}
		x[key] = value
	}
	return x, nil
}

// nextRegularToken return the next non-whitespace token
func (p *Parser) nextRegularToken() token {
	token := p.nextToken()
	if token.typ == whitespaceType {
		return p.nextToken()
	}
	return token
}

// nextToken returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) nextToken() token {
	// If we have a token on the buffer, then return it.
	if !p.ts.isEmpty() {
		token, err := p.ts.pop()
		if err != nil {
			panic(err)
		}
		return token
	}
	// Otherwise read the next token from the scanner.
	return p.lex.Lex()
}

// backup pushes the a token back onto the stack.
func (p *Parser) backup(tok token) {
	p.ts.push(tok)
}

func (p *Parser) makeError(format string, a ...any) error {
	s := fmt.Sprintf(format, a...)
	return fmt.Errorf("%s in line %d", s, p.lex.loc)
}
