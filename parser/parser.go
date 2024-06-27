package parser

import (
	"fmt"
	"io"
	"strconv"
)

type tokenBuffer struct {
	tok   Token
	value any
	n     int
}

// Parser represents a parser.
type Parser struct {
	l   *Lexer
	buf tokenBuffer // last read token
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{l: NewScanner(r)}
}
func (p *Parser) Parse() (map[string]any, error) {
	x := make(map[string]any)

loop:
	for {
		var key string
		var value any

		// First token should be identifier or integer
		switch tok, v := p.scanIgnoreWhitespace(); tok {
		case Eof, BracketsClose:
			break loop
		case Identifier:
			key = v.(string)
		case Integer:
			key = strconv.Itoa(v.(int))
		default:
			return nil, fmt.Errorf("found %v, expected identifier or integer", v)
		}

		// Next should be an equal sign
		if tok, lit := p.scanIgnoreWhitespace(); tok != Equal {
			return nil, fmt.Errorf("found %v, expected equal sign", lit)
		}

		// Next should be some kind of value
		switch tok, v := p.scanIgnoreWhitespace(); tok {
		case String, Float, Integer, Boolean:
			value = v
		case BracketsOpen:
			tok, _ := p.scanIgnoreWhitespace()
			p.unscan()
			switch tok {
			case Identifier:
				v, err := p.Parse()
				if err != nil {
					return nil, err
				}
				value = v
			case Integer:
				x := make([]int, 0)
				for {
					tok, v := p.scanIgnoreWhitespace()
					if tok == BracketsClose {
						value = x
						break
					}
					x = append(x, v.(int))
				}
			case String:
				x := make([]string, 0)
				for {
					tok, v := p.scanIgnoreWhitespace()
					if tok == BracketsClose {
						value = x
						break
					}
					x = append(x, v.(string))
				}
			default:
				return nil, fmt.Errorf("invalid token %v for array", v)
			}

		default:
			return nil, fmt.Errorf("found %v, expected a value", v)
		}
		x[key] = value
	}
	return x, nil
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, value any) {
	tok, value = p.scan()
	if tok == Whitespace {
		tok, value = p.scan()
	}
	return
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, value any) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.value
	}

	// Otherwise read the next token from the scanner.
	tok, value = p.l.Lex()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.value = tok, value

	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() {
	p.buf.n = 1
}
