package parser

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"unicode"
)

var eof = rune(0)

// Lexer represents a lexical scanner.
type Lexer struct {
	r   *bufio.Reader
	loc int
}

// NewLexer returns a new instance of lexer
func NewLexer(r io.Reader) *Lexer {
	return &Lexer{r: bufio.NewReader(r), loc: 1}
}

// Lex returns the next token and literal value.
func (l *Lexer) Lex() token {
	// Read the next rune.
	for {
		ch := l.read()
		if unicode.IsSpace(ch) {
			l.unread()
			l.consumeWhitespace()
			continue
		}
		if ch == '"' {
			return l.scanString()
		}
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '-' {
			l.unread()
			return l.scanWord()
		}
		switch ch {
		case eof:
			return token{endOfFile, ""}
		case '{':
			return token{bracketsOpen, string(ch)}
		case '}':
			return token{bracketsClose, string(ch)}
		case '=':
			return token{equalSign, string(ch)}
		}
		return token{illegal, string(ch)}
	}
}

// read reads the next rune from the buffered reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (l *Lexer) read() rune {
	ch, _, err := l.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (l *Lexer) unread() {
	_ = l.r.UnreadRune()
}

// consumeWhitespace consumes all whitespace from the reader.
func (l *Lexer) consumeWhitespace() {
	for {
		ch := l.read()
		if ch == eof {
			break
		}
		if !unicode.IsSpace(ch) {
			l.unread()
			break
		}
		if ch == '\n' {
			l.loc++
		}
	}
}

// scanWord identifiers a word, which can be an identifier, a keyword or a number.
func (l *Lexer) scanWord() token {
	var buf bytes.Buffer
	buf.WriteRune(l.read())

	// Read word from stream
	for {
		if ch := l.read(); ch == eof {
			break
		} else if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' && ch != '-' && ch != '.' {
			l.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	s := buf.String()
	hasLetter := false
	for _, x := range s {
		if unicode.IsLetter(x) || x == '_' {
			hasLetter = true
			break
		}
	}
	if !hasLetter {
		x1, err := strconv.ParseFloat(s, 64)
		if err == nil { // if this was no float after we assume a string
			x2 := int(x1)
			if x1 == float64(x2) {
				return token{integer, x2}
			}
			return token{float, x1}
		}
	}

	// If the word matches a keyword then return that that token.
	switch s {
	case "yes":
		return token{boolean, true}
	case "no":
		return token{boolean, false}
	}
	// Otherwise return as a identifier.
	return token{identifier, s}
}

// scanWord identifies a string token.
func (l *Lexer) scanString() token {
	var buf bytes.Buffer
	for {
		ch := l.read()
		if ch == eof || ch == '"' {
			break
		}
		_, _ = buf.WriteRune(ch)
	}
	s := buf.String()
	return token{str, s}
}
