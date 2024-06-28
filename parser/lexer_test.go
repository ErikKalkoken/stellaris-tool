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
		{"-42", token{integerType, -42}},
		{"{", token{bracketsOpenType, "{"}},
		{"}", token{bracketsCloseType, "}"}},
		{" ", token{eofType, ""}},
		{" 			 ", token{eofType, ""}},
		{"#", token{illegalType, "#"}},
		// special words
		{"yes", token{booleanType, true}},
		{"no", token{booleanType, false}},
		{"none", token{identifierType, "none"}},
		{"not_set", token{identifierType, "not_set"}},
		{"indeterminable", token{identifierType, "indeterminable"}},
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
			[]tokenType{identifierType, identifierType},
		},
		{
			"hello    	   world",
			[]tokenType{identifierType, identifierType},
		},
		{
			"yes no hello",
			[]tokenType{booleanType, booleanType, identifierType},
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
			"first=123.45 second=5",
			[]tokenType{identifierType, equalSignType, floatType, identifierType, equalSignType, integerType},
		},
		{
			"first=none",
			[]tokenType{identifierType, equalSignType, identifierType},
		},
		{
			"x={next_usable_date=\"-5070.07.21\"}",
			[]tokenType{identifierType, equalSignType, bracketsOpenType, identifierType, equalSignType, stringType, bracketsCloseType},
		},
		{
			"\"\" first=\"\"",
			[]tokenType{stringType, identifierType, equalSignType, stringType},
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

func TestSpecialFeatures(t *testing.T) {
	t.Run("can keep track of LOC", func(t *testing.T) {
		in := strings.NewReader("alpha=1\nbravo=2")
		s := NewLexer(in)
		for {
			token := s.Lex()
			if token.typ == eofType {
				break
			}
		}
		assert.Equal(t, 2, s.loc)
	})
}
