package parser

import (
	"bufio"
	"bytes"
	"io"
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
	Number
	Yes
	No
)

var eof = rune(0)

// Scanner represents a lexical scanner.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (tok Token, lit string) {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	if unicode.IsSpace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if ch == '"' {
		return s.scanString()
	} else if unicode.IsLetter(ch) {
		s.unread()
		return s.scanWord()
	} else if unicode.IsDigit(ch) {
		s.unread()
		return s.scanNumber()
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
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !unicode.IsSpace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return Whitespace, buf.String()
}

// scanWord consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanWord() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}
	// If the string matches a keyword then return that keyword.
	v := buf.String()
	switch v {
	case "yes":
		return Yes, v
	case "no":
		return No, v
	}

	// Otherwise return as a regular identifier.
	return Identifier, v
}

// scanWord consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanString() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
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
func (s *Scanner) scanNumber() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !unicode.IsDigit(ch) && ch != '.' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return Number, buf.String()
}
