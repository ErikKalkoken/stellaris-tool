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
		{"name", token{identifierType, "name"}},
		{"\"string\"", token{stringType, "string"}},
		{"1.234", token{floatType, 1.234}},
		{"42", token{integerType, 42}},
		{"{", token{bracketsOpenType, "{"}},
		{"}", token{bracketsCloseType, "}"}},
		{" ", token{whitespaceType, " "}},
		{" 			 ", token{whitespaceType, " 			 "}},
		{"#", token{illegalType, "#"}},
		// keywords
		{"yes", token{booleanType, true}},
		{"no", token{booleanType, false}},
		{"not_set", token{keywordType, NotSet}},
		{"indeterminable", token{keywordType, Indeterminable}},
		{"none", token{keywordType, None}},
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
			[]tokenType{identifierType, whitespaceType, identifierType},
		},
		{
			"hello    	   world",
			[]tokenType{identifierType, whitespaceType, identifierType},
		},
		{
			"yes no hello",
			[]tokenType{booleanType, whitespaceType, booleanType, whitespaceType, identifierType},
		},
		{
			"first=\"second 123 $%&\"",
			[]tokenType{identifierType, equalSignType, stringType},
		},
		{
			"first=123.45",
			[]tokenType{identifierType, equalSignType, floatType},
		},
		{
			"first=123.45second=5",
			[]tokenType{identifierType, equalSignType, floatType, identifierType, equalSignType, integerType},
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("in: %s", tc.in), func(t *testing.T) {
			in := strings.NewReader(tc.in)
			s := NewLexer(in)
			got := make([]tokenType, 0)
			for {
				token := s.Lex()
				if token.typ == eofType {
					break
				}
				got = append(got, token.typ)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
