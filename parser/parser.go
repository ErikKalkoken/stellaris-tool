// Package parser contains a parser for Paradox save files.
package parser

import (
	"fmt"
	"io"
	"strconv"
)

var emptyObject = struct{}{}

// Parser represents a parser for Paradox save files.
type Parser struct {
	// Provides a stream of tokens
	lex *lexer
	// Stack of latest tokens so we can go back
	ts stack[token]
}

// NewParser takes a reader and returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{lex: newLexer(r), ts: newStack[token](3)}
}

// Parse parsed a Paradox save file and returns it's contents.
//
// Here is how the parser deals with some particulars of the paradox format:
// - The format allows multiple values for a key, so the parser returns a nested map of keys to value slices.
// - All keys are converted to strings, including keywords and numbers
// - The keywords "none" and "not_set" are converted to nil (when used as values)
// - The keywords "yes" and "no" are converted to bool
// - An array of strings or boolean will be returned as string or bool slices respectively
// - An array of numbers will be returns as a float64 slice
// - Arrays can also be empty
func (p *Parser) Parse() (map[string][]any, error) {
	result := make(map[string][]any)
loop:
	for {
		var key string
		var value any

		// First token should some kind of key or signaling the end of the current nesting level
		switch tok := p.nextToken(); tok.typ {
		case endOfFile, bracketsClose:
			break loop
		case identifier, str:
			key = tok.value.(string)
		case integer:
			key = strconv.Itoa(tok.value.(int))
		default:
			return nil, p.makeError("found %v, expected some kind of key", tok)
		}

		// Next is usually an equal sign. If it is omitted we assume there is one.
		if tok := p.nextToken(); tok.typ != equalSign {
			p.backup(tok)
		}

		// Next should be some kind of value
		switch tok := p.nextToken(); tok.typ {
		case str, float, integer, boolean:
			value = tok.value
		case identifier:
			if tok.value == "none" || tok.value == "not_set" {
				value = nil
			} else {
				value = tok.value
			}
		case bracketsOpen:
			tok2 := p.nextToken()
			switch tok2.typ {
			case bracketsClose:
				// Empty object
				value = emptyObject
			case bracketsOpen:
				// Array of objects
				oo := make([]map[string][]any, 0)
				for {
					v2, err := p.Parse()
					if err != nil {
						return nil, err
					}
					oo = append(oo, v2)
					tok3 := p.nextToken()
					if tok3.typ == bracketsClose {
						break
					} else if tok3.typ != bracketsOpen {
						return nil, p.makeError("unexpected token %v in obj array", tok3)
					}
				}
				value = oo
			case identifier, str:
				tok3 := p.nextToken()
				if tok3.typ == equalSign {
					// A regular object
					p.backup(tok3)
					p.backup(tok2)
					x, err := p.Parse()
					if err != nil {
						return nil, err
					}
					value = x
				} else {
					// Array of string
					p.backup(tok3)
					p.backup(tok2)
					ss := make([]string, 0)
					for {
						tok3 := p.nextToken()
						if tok3.typ == bracketsClose {
							break
						}
						y, ok := tok3.value.(string)
						if !ok {
							return nil, p.makeError("found %v, expected type string for array", tok3)
						}
						ss = append(ss, y)
					}
					value = ss
				}
			case integer, float:
				if tok2.typ == integer {
					tok3 := p.nextToken()
					p.backup(tok3)
					p.backup(tok2)
					if tok3.typ == equalSign {
						// An ID object
						x, err := p.Parse()
						if err != nil {
							return nil, err
						}
						value = x
						break
					}
					if tok3.typ == bracketsOpen {
						panic(p.makeError("unexpected token: %v", tok3))
					}
				} else {
					p.backup(tok2)
				}
				// Array of numbers
				ff := make([]float64, 0)
				for {
					tok3 := p.nextToken()
					if tok3.typ == bracketsClose {
						break
					}
					switch tok3.typ {
					case float:
						ff = append(ff, tok3.value.(float64))
					case integer:
						ff = append(ff, float64(tok3.value.(int)))
					default:
						return nil, p.makeError("unexpected token for number array: %v", tok3)
					}
				}
				value = ff
			case boolean:
				// Array of boolean
				p.backup(tok2)
				ss := make([]bool, 0)
				for {
					tok3 := p.nextToken()
					if tok3.typ == bracketsClose {
						break
					}
					y, ok := tok3.value.(bool)
					if !ok {
						return nil, p.makeError("expected type boolean for boolean array, but got: %v", tok3)
					}
					ss = append(ss, y)
				}
				value = ss
			default:
				return nil, p.makeError("invalid token %v for array", tok2)
			}

		default:
			return nil, p.makeError("found %v, expected a value", tok)
		}
		if value != emptyObject {
			result[key] = append(result[key], value)
		} else {
			result[key] = make([]any, 0)
		}
	}
	return result, nil
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
	token := p.lex.lex()
	// fmt.Println(token)
	return token
}

// backup pushes the a token back onto the stack.
func (p *Parser) backup(tok token) {
	p.ts.push(tok)
}

func (p *Parser) makeError(format string, a ...any) error {
	s := fmt.Sprintf(format, a...)
	return fmt.Errorf("%s in line %d", s, p.lex.loc)
}
