package parser

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"unicode"
)

type Token uint

const (
	Illegal Token = iota
	Eof
	Whitespace
	Equal
	BracketsOpen
	BracketsClose
	Identifier
	String
	Float
	Integer
	Boolean
)

var eof = rune(0)

// Lexer represents a lexical scanner.
type Lexer struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Lexer {
	return &Lexer{r: bufio.NewReader(r)}
}

// Lex returns the next token and literal value.
func (l *Lexer) Lex() (tok Token, value any) {
	// Read the next rune.
	ch := l.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	if unicode.IsSpace(ch) {
		l.unread()
		return l.scanWhitespace()
	} else if ch == '"' {
		return l.scanString()
	} else if unicode.IsLetter(ch) {
		l.unread()
		return l.scanWord()
	} else if unicode.IsDigit(ch) {
		l.unread()
		return l.scanNumber()
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return Eof, ""
	case '{':
		return BracketsOpen, string(ch)
	case '}':
		return BracketsClose, string(ch)
	case '=':
		return Equal, string(ch)
	}

	return Illegal, string(ch)
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

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (l *Lexer) scanWhitespace() (tok Token, value any) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(l.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := l.read(); ch == eof {
			break
		} else if !unicode.IsSpace(ch) {
			l.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return Whitespace, buf.String()
}

// scanWord consumes the current rune and all contiguous ident runes.
func (l *Lexer) scanWord() (tok Token, value any) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(l.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := l.read(); ch == eof {
			break
		} else if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
			l.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}
	// If the string matches a keyword then return that keyword.
	s := buf.String()
	switch s {
	case "yes":
		return Boolean, true
	case "no":
		return Boolean, false
	}

	// Otherwise return as a regular identifier.
	return Identifier, s
}

// scanWord consumes the current rune and all contiguous ident runes.
func (l *Lexer) scanString() (tok Token, value string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(l.read())

	for {
		if ch := l.read(); ch == eof {
			break
		} else if ch == '"' {
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}
	return String, buf.String()
}

// scanWord consumes the current rune and all contiguous ident runes.
func (l *Lexer) scanNumber() (tok Token, value any) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(l.read())

	for {
		if ch := l.read(); ch == eof {
			break
		} else if !unicode.IsDigit(ch) && ch != '.' {
			l.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	s := buf.String()
	x1, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	x2 := int(x1)
	if x1 == float64(x2) {
		return Integer, x2
	}
	return Float, x1
}
