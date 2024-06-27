package parser

import (
	"fmt"
	"io"
	"strconv"
)

type tokenBuffer struct {
	tok Token  // last read token
	lit string // last read literal
	n   int    // buffer size (max=1)
}

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf tokenBuffer
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func (p *Parser) Parse() (map[string]any, error) {
	var key string
	var value any
	// First token should be an identifier
	if tok, lit := p.scanIgnoreWhitespace(); tok != Identifier {
		return nil, fmt.Errorf("found %q, expected identifier", lit)
	} else {
		key = lit
	}

	// Next should be an equal sign
	if tok, lit := p.scanIgnoreWhitespace(); tok != Equal {
		return nil, fmt.Errorf("found %q, expected equal sign", lit)
	}

	// Next should be a value
	tok, lit := p.scanIgnoreWhitespace()
	switch tok {
	case String:
		value = lit
	case Number:
		x, err := strconv.ParseFloat(lit, 64)
		if err != nil {
			return nil, err
		}
		y := int(x)
		if x == float64(y) {
			value = y
		} else {
			value = x
		}
	default:
		return nil, fmt.Errorf("found %q, expected a value", lit)
	}
	return map[string]any{key: value}, nil
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == Whitespace {
		tok, lit = p.scan()
	}
	return
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() {
	p.buf.n = 1
}
