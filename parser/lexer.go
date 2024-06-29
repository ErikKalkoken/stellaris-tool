package parser

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"unicode"
)

var eof = rune(0)

// lexer represents a lexical scanner.
type lexer struct {
	r   *bufio.Reader
	loc int
}

// newLexer returns a new instance of lexer
func newLexer(r io.Reader) *lexer {
	return &lexer{r: bufio.NewReader(r), loc: 1}
}

// lex returns the next token and literal value. This is the main method.
func (l *lexer) lex() (token, error) {
	// Read the next rune.
	for {
		ch, err := l.read()
		if err != nil {
			return token{}, err
		}
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
			return token{endOfFile, ""}, nil
		case '{':
			return token{bracketsOpen, string(ch)}, nil
		case '}':
			return token{bracketsClose, string(ch)}, nil
		case '=':
			return token{equalSign, string(ch)}, nil
		}
		return token{illegal, string(ch)}, nil
	}
}

// read reads and returns the next rune from the buffered reader or the EOF rune.
func (l *lexer) read() (rune, error) {
	ch, _, err := l.r.ReadRune()
	if err == io.EOF {
		return eof, nil
	} else if err != nil {
		return 0, err
	}
	return ch, nil
}

// unread places the previously read rune back on the reader.
func (l *lexer) unread() error {
	if err := l.r.UnreadRune(); err != nil {
		return err
	}
	return nil
}

// consumeWhitespace consumes all whitespace from the reader.
func (l *lexer) consumeWhitespace() error {
	for {
		ch, err := l.read()
		if err != nil {
			return err
		}
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
	return nil
}

// scanWord returns an identifier, a keyword or a number from the scanned input.
func (l *lexer) scanWord() (token, error) {
	var buf bytes.Buffer
	ch, err := l.read()
	if err != nil {
		return token{}, err
	}
	_, err = buf.WriteRune(ch)
	if err != nil {
		return token{}, err
	}
	// Read word from stream into a buffer
	for {
		ch, err := l.read()
		if err != nil {
			return token{}, err
		}
		if ch == eof {
			break
		} else if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' && ch != '-' && ch != '.' {
			l.unread()
			break
		} else {
			_, err := buf.WriteRune(ch)
			if err != nil {
				return token{}, err
			}
		}
	}
	// parse the string
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
		if err == nil { // if this was actually no float we assume it's a string (e.g. could be a date)
			x2 := int(x1)
			if x1 == float64(x2) {
				return token{integer, x2}, nil
			}
			return token{float, x1}, nil
		}
	}

	// If the word matches a keyword then return that that token.
	switch s {
	case "yes":
		return token{boolean, true}, nil
	case "no":
		return token{boolean, false}, nil
	}
	// Otherwise return as a identifier.
	return token{identifier, s}, nil
}

// scanString returns a string token from the scanned input.
func (l *lexer) scanString() (token, error) {
	var buf bytes.Buffer
	for {
		ch, err := l.read()
		if err != nil {
			return token{}, err
		}
		if ch == eof || ch == '"' {
			break
		}
		_, err = buf.WriteRune(ch)
		if err != nil {
			return token{}, err
		}
	}
	s := buf.String()
	return token{str, s}, nil
}
