package parser

import (
	"fmt"
	"io"
	"strconv"
)

// Parser represents a parser.
type Parser struct {
	// Provides a stream of tokens
	lex *Lexer
	// Stack of latest tokens so we can go back
	ts stack[Token]
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{lex: NewLexer(r), ts: newStack[Token](3)}
}
func (p *Parser) Parse() (map[string]any, error) {
	x := make(map[string]any)

loop:
	for {
		var key string
		var value any

		// First token should be identifier or integer
		switch tok := p.nextRegularToken(); tok.typ {
		case Eof, BracketsClose:
			break loop
		case Identifier:
			key = tok.value.(string)
		case Integer:
			key = strconv.Itoa(tok.value.(int))
		default:
			return nil, p.makeError("found %v, expected identifier or integer", tok)
		}

		// Next should be an equal sign
		if tok := p.nextRegularToken(); tok.typ != EqualSign {
			return nil, p.makeError("found %v, expected equal sign", tok)
		}

		// Next should be some kind of value
		switch tok := p.nextRegularToken(); tok.typ {
		case String, Float, Integer, Boolean:
			value = tok.value
		case BracketsOpen:
			tok2 := p.nextRegularToken()
			switch tok2.typ {
			case BracketsClose:
				// Empty object
				value = struct{}{}
			case BracketsOpen:
				// Array of objects
				x := make([]map[string]any, 0)
				for {
					v2, err := p.Parse()
					if err != nil {
						return nil, err
					}
					x = append(x, v2)
					tok3 := p.nextRegularToken()
					if tok3.typ == BracketsClose {
						break
					} else if tok3.typ != BracketsOpen {
						return nil, p.makeError("Unexpected token %v in obj array", tok3)
					}
				}
				value = x
			case Identifier:
				// Normal object
				p.backup(tok2)
				x, err := p.Parse()
				if err != nil {
					return nil, err
				}
				value = x
			case Integer:
				tok3 := p.nextRegularToken()
				p.backup(tok3)
				p.backup(tok2)
				if tok3.typ == EqualSign {
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
					if tok3.typ == BracketsClose {
						value = x
						break
					}
					y, ok := tok3.value.(int)
					if !ok {
						return nil, p.makeError("Expected type integer, but got: %v", tok3)
					}
					x = append(x, y)
				}
			case Float:
				// Array of float
				p.backup(tok2)
				x := make([]float64, 0)
				for {
					tok3 := p.nextRegularToken()
					if tok3.typ == BracketsClose {
						value = x
						break
					}
					y, ok := tok3.value.(float64)
					if !ok {
						return nil, p.makeError("Expected type float, but got: %v", tok2)
					}
					x = append(x, y)
				}
			case String:
				// Array of string
				p.backup(tok2)
				x := make([]string, 0)
				for {
					tok3 := p.nextRegularToken()
					if tok3.typ == BracketsClose {
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
func (p *Parser) nextRegularToken() Token {
	token := p.nextToken()
	if token.typ == Whitespace {
		return p.nextToken()
	}
	return token
}

// nextToken returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) nextToken() Token {
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
func (p *Parser) backup(tok Token) {
	p.ts.push(tok)
}

func (p *Parser) makeError(format string, a ...any) error {
	s := fmt.Sprintf(format, a...)
	return fmt.Errorf("%s in line %d", s, p.lex.loc)
}
