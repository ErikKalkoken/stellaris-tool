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
		want token
	}{
		{"name", token{identifier, "name"}},
		{"\"string\"", token{str, "string"}},
		{"1.234", token{float, 1.234}},
		{"42", token{integer, 42}},
		{"-42", token{integer, -42}},
		{"{", token{bracketsOpen, "{"}},
		{"}", token{bracketsClose, "}"}},
		{" ", token{endOfFile, ""}},
		{" 			 ", token{endOfFile, ""}},
		{"#", token{illegal, "#"}},
		// special words
		{"yes", token{boolean, true}},
		{"no", token{boolean, false}},
		{"none", token{identifier, "none"}},
		{"not_set", token{identifier, "not_set"}},
		{"indeterminable", token{identifier, "indeterminable"}},
		{`"one \"two\" three"`, token{str, "one \"two\" three"}},
		{`"one \\ two"`, token{str, "one \\ two"}},
		{"one:two", token{identifier, "one:two"}},
		{"@one", token{identifier, "@one"}},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("in: %s", tc.in), func(t *testing.T) {
			in := strings.NewReader(tc.in)
			l := newLexer(in)
			got, _ := l.lex()
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
			[]tokenType{identifier, identifier},
		},
		{
			"hello    	   world",
			[]tokenType{identifier, identifier},
		},
		{
			"yes no hello",
			[]tokenType{boolean, boolean, identifier},
		},
		{
			"first=\"second 123 $%&\"",
			[]tokenType{identifier, equalSign, str},
		},
		{
			"first=123.45",
			[]tokenType{identifier, equalSign, float},
		},
		{
			"first=123.45 second=5",
			[]tokenType{identifier, equalSign, float, identifier, equalSign, integer},
		},
		{
			"first=none",
			[]tokenType{identifier, equalSign, identifier},
		},
		{
			"x={next_usable_date=\"-5070.07.21\"}",
			[]tokenType{identifier, equalSign, bracketsOpen, identifier, equalSign, str, bracketsClose},
		},
		{
			"\"\" first=\"\"",
			[]tokenType{str, identifier, equalSign, str},
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("in: %s", tc.in), func(t *testing.T) {
			in := strings.NewReader(tc.in)
			s := newLexer(in)
			got := make([]tokenType, 0)
			for {
				token, _ := s.lex()
				if token.typ == endOfFile {
					break
				}
				got = append(got, token.typ)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestSpecialFeatures(t *testing.T) {
	t.Run("can keep track of LOC", func(t *testing.T) {
		in := strings.NewReader("alpha=1\nbravo=2")
		s := newLexer(in)
		for {
			token, _ := s.lex()
			if token.typ == endOfFile {
				break
			}
		}
		assert.Equal(t, 2, s.loc)
	})
}
