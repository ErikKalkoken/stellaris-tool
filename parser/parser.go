package parser

import (
	"fmt"
	"io"
	"strconv"
)

type tokenBuffer struct {
	token Token
	n     int
}

// Parser represents a parser.
type Parser struct {
	l   *Lexer
	buf tokenBuffer // last read token
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{l: NewLexer(r)}
}
func (p *Parser) Parse() (map[string]any, error) {
	x := make(map[string]any)

loop:
	for {
		var key string
		var value any

		// First token should be identifier or integer
		switch tok := p.scanIgnoreWhitespace(); tok.typ {
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
		if tok := p.scanIgnoreWhitespace(); tok.typ != Equal {
			return nil, p.makeError("found %v, expected equal sign", tok)
		}

		// Next should be some kind of value
		switch tok := p.scanIgnoreWhitespace(); tok.typ {
		case String, Float, Integer, Boolean:
			value = tok.value
		case BracketsOpen:
			tok2 := p.scanIgnoreWhitespace()
			if tok2.typ == BracketsOpen {
				// Array of objects
				x := make([]map[string]any, 0)
				for {
					v2, err := p.Parse()
					if err != nil {
						return nil, err
					}
					x = append(x, v2)
					tok3 := p.scanIgnoreWhitespace()
					if tok3.typ == BracketsClose {
						value = x
						break
					} else if tok3.typ != BracketsOpen {
						return nil, p.makeError("Unexpected token %v in obj array", tok3)
					}
				}
			} else {
				// Array of value
				switch tok2.typ {
				case Identifier:
					p.unscan()
					x, err := p.Parse()
					if err != nil {
						return nil, err
					}
					value = x
				case Integer:

					p.unscan()
					x := make([]int, 0)
					for {
						tok3 := p.scanIgnoreWhitespace()
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
					p.unscan()
					x := make([]float64, 0)
					for {
						tok3 := p.scanIgnoreWhitespace()
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
					p.unscan()
					x := make([]string, 0)
					for {
						tok3 := p.scanIgnoreWhitespace()
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
			}

		default:
			return nil, p.makeError("found %v, expected a value", tok)
		}
		x[key] = value
	}
	return x, nil
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() Token {
	token := p.scan()
	if token.typ == Whitespace {
		return p.scan()
	}
	return token
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() Token {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.token
	}

	// Otherwise read the next token from the scanner.
	token := p.l.Lex()

	// Save it to the buffer in case we unscan later.
	p.buf.token = token

	return token
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() {
	p.buf.n = 1
}

func (p *Parser) makeError(format string, a ...any) error {
	s := fmt.Sprintf(format, a...)
	return fmt.Errorf("%s in line %d", s, p.l.loc)
}
