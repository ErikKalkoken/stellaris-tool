package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleTokens(t *testing.T) {
	cases := []struct {
		in   string
		want Token
	}{
		{"name", Token{Identifier, "name"}},
		{"\"string\"", Token{String, "string"}},
		{"1.234", Token{Float, 1.234}},
		{"yes", Token{Boolean, true}},
		{"no", Token{Boolean, false}},
		{"{", Token{BracketsOpen, "{"}},
		{"}", Token{BracketsClose, "}"}},
		{" ", Token{Whitespace, " "}},
		{" 			 ", Token{Whitespace, " 			 "}},
		{"#", Token{Illegal, "#"}},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("in: %s", tc.in), func(t *testing.T) {
			in := strings.NewReader(tc.in)
			l := NewLexer(in)
			got := l.Lex()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestMultipleTokens(t *testing.T) {
	cases := []struct {
		in   string
		want []tokenType
	}{
		{
			"hello world",
			[]tokenType{Identifier, Whitespace, Identifier},
		},
		{
			"hello    	   world",
			[]tokenType{Identifier, Whitespace, Identifier},
		},
		{
			"yes no hello",
			[]tokenType{Boolean, Whitespace, Boolean, Whitespace, Identifier},
		},
		{
			"first=\"second 123 $%&\"",
			[]tokenType{Identifier, Equal, String},
		},
		{
			"first=123.45",
			[]tokenType{Identifier, Equal, Float},
		},
		{
			"first=123.45second=5",
			[]tokenType{Identifier, Equal, Float, Identifier, Equal, Integer},
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("in: %s", tc.in), func(t *testing.T) {
			in := strings.NewReader(tc.in)
			s := NewLexer(in)
			got := make([]tokenType, 0)
			for {
				token := s.Lex()
				if token.typ == Eof {
					break
				}
				got = append(got, token.typ)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
