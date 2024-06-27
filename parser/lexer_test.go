package parser_test

import (
	"example/stellaris-tool/parser"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type result struct {
	token parser.Token
	value any
}

func TestSingleTokens(t *testing.T) {
	cases := []struct {
		in   string
		want result
	}{
		{"name", result{parser.Identifier, "name"}},
		{"\"string\"", result{parser.String, "string"}},
		{"1.234", result{parser.Float, 1.234}},
		{"yes", result{parser.Boolean, true}},
		{"no", result{parser.Boolean, false}},
		{"{", result{parser.BracketsOpen, "{"}},
		{"}", result{parser.BracketsClose, "}"}},
		{" ", result{parser.Whitespace, " "}},
		{" 			 ", result{parser.Whitespace, " 			 "}},
		{"#", result{parser.Illegal, "#"}},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("in: %s", tc.in), func(t *testing.T) {
			in := strings.NewReader(tc.in)
			l := parser.NewLexer(in)
			token, lit := l.Lex()
			got := result{token: token, value: lit}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestMultipleTokens(t *testing.T) {
	cases := []struct {
		in   string
		want []parser.Token
	}{
		{
			"hello world",
			[]parser.Token{parser.Identifier, parser.Whitespace, parser.Identifier},
		},
		{
			"hello    	   world",
			[]parser.Token{parser.Identifier, parser.Whitespace, parser.Identifier},
		},
		{
			"yes no hello",
			[]parser.Token{parser.Boolean, parser.Whitespace, parser.Boolean, parser.Whitespace, parser.Identifier},
		},
		{
			"first=\"second 123 $%&\"",
			[]parser.Token{parser.Identifier, parser.Equal, parser.String},
		},
		{
			"first=123.45",
			[]parser.Token{parser.Identifier, parser.Equal, parser.Float},
		},
		{
			"first=123.45second=5",
			[]parser.Token{parser.Identifier, parser.Equal, parser.Float, parser.Identifier, parser.Equal, parser.Integer},
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("in: %s", tc.in), func(t *testing.T) {
			in := strings.NewReader(tc.in)
			s := parser.NewLexer(in)
			got := make([]parser.Token, 0)
			for {
				token, _ := s.Lex()
				if token == parser.Eof {
					break
				}
				got = append(got, token)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
