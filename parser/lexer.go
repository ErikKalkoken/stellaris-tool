package parser

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"strings"
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
	ch := l.read()

	if unicode.IsSpace(ch) {
		l.unread()
		return l.scanWhitespace()
	} else if ch == '"' {
		return l.scanString()
	} else if unicode.IsLetter(ch) {
		l.unread()
		return l.scanWord()
	} else if unicode.IsDigit(ch) || ch == '-' {
		l.unread()
		return l.scanNumber()
	}

	switch ch {
	case eof:
		return token{eofType, ""}
	case '{':
		return token{bracketsOpenType, string(ch)}
	case '}':
		return token{bracketsCloseType, string(ch)}
	case '=':
		return token{equalSignType, string(ch)}
	}

	return token{illegalType, string(ch)}
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

// scanWhitespace identifies a whitespace token, which can contain multiple continuous whitespace runes.
// It also updates the current line number.
func (l *Lexer) scanWhitespace() token {
	var buf bytes.Buffer
	buf.WriteRune(l.read())

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

	s := buf.String()
	l.loc += strings.Count(s, "\n")
	return token{whitespaceType, s}
}

// scanWord identifiers a word, which can be an identifier or a keyword.
func (l *Lexer) scanWord() token {
	var buf bytes.Buffer
	buf.WriteRune(l.read())

	// Read word from stream
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
	// If the word matches a keyword then return that that token.
	s := buf.String()
	switch s {
	case "yes":
		return token{booleanType, true}
	case "no":
		return token{booleanType, false}
	case "not_set":
		return token{keywordType, NotSet}
	case "indeterminable":
		return token{keywordType, Indeterminable}
	case "male":
		return token{keywordType, Male}
	case "female":
		return token{keywordType, Female}
	case "none":
		return token{keywordType, None}
	}
	// Otherwise return as a identifier.
	return token{identifierType, s}
}

// scanWord identifies a string token.
func (l *Lexer) scanString() token {
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
	return token{stringType, buf.String()}
}

// scanNumber identifiers a number type token.
func (l *Lexer) scanNumber() token {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(l.read())

	for {
		if ch := l.read(); ch == eof {
			break
		} else if !unicode.IsDigit(ch) && ch != '.' && ch != '-' {
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
		return token{integerType, x2}
	}
	return token{floatType, x1}
}
