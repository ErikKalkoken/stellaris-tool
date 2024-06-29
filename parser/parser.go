package parser

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	duplicateKeyPrefix = "DUPLICATE_KEY_"
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
	duplicateKeys := map[string]int{}
loop:
	for {
		var key string
		var value any

		// First token should be identifier or integer or string
		switch tok := p.nextToken(); tok.typ {
		case endOfFile, bracketsClose:
			break loop
		case identifier, str:
			key = tok.value.(string)
		case integer:
			key = strconv.Itoa(tok.value.(int))
		default:
			return nil, p.makeError("found %v, expected identifier or integer", tok)
		}

		// Next is usually an equal sign, or we assume one
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
				value = struct{}{}
			case bracketsOpen:
				// Array of objects
				oo := make([]map[string]any, 0)
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
						return nil, p.makeError("Unexpected token %v in obj array", tok3)
					}
				}
				value = oo
			case identifier, str:
				tok3 := p.nextToken()
				if tok3.typ == equalSign {
					// object
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
							return nil, p.makeError("Expected type string for array, but got: %v", tok3)
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
						// ID object
						x, err := p.Parse()
						if err != nil {
							return nil, err
						}
						value = x
						break
					}
					if tok3.typ == bracketsOpen {
						panic(p.makeError("Unexpected token: %v", tok3))
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
						return nil, p.makeError("Unexpected token for float array: %v", tok3)
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
						return nil, p.makeError("Expected type boolean for array, but got: %v", tok3)
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
		// handle duplicate keys part 1 (when a new k/v pair is added)
		_, found := x[key]
		if found {
			duplicateKeys[key] = 0
			newKey := makeKeyWithSuffix(key, 0)
			x[newKey] = x[key]
			delete(x, key)
		}
		_, found = duplicateKeys[key]
		if found {
			duplicateKeys[key]++
			key = makeKeyWithSuffix(key, duplicateKeys[key])
		}
		// handle duplicate keys part 2 (when we have all duplicate k/v pairs)
		m, ok := value.(map[string]any)
		if ok {
			duplicates := make(map[string]map[int]any)
			for k, v := range m {
				k2, found := strings.CutPrefix(k, duplicateKeyPrefix)
				if found {
					p := strings.SplitN(k2, "_", 2)
					if len(p) != 2 {
						return nil, fmt.Errorf("duplicate key has unexpected format: %s", k2)
					}
					k3 := p[1]
					id, err := strconv.Atoi(p[0])
					if err != nil {
						return nil, fmt.Errorf("failed to convert duplicate key ID: %s", k2)
					}
					_, found := duplicates[k3]
					if !found {
						duplicates[k3] = make(map[int]any)
					}
					duplicates[k3][id] = v
					delete(m, k)
				}
			}
			for k, m2 := range duplicates {
				a := make([]any, len(m2))
				for id, v := range m2 {
					a[id] = v
				}
				m[k] = a
			}
		}
		x[key] = value
	}
	return x, nil
}

func makeKeyWithSuffix(key string, id int) string {
	return fmt.Sprintf("%s%d_%s", duplicateKeyPrefix, id, key)
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
	token := p.lex.Lex()
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
