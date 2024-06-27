package parser_test

import (
	"example/stellaris-tool/parser"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	cases := []struct {
		in   string
		want map[string]any
	}{
		{"alpha=5", map[string]any{"alpha": 5}},
		{"alpha=5.3", map[string]any{"alpha": 5.3}},
		{"alpha=\"special text\"", map[string]any{"alpha": "special text"}},
	}
	for _, tc := range cases {
		r := strings.NewReader(tc.in)
		p := parser.NewParser(r)
		got, err := p.Parse()
		if assert.NoError(t, err) {
			assert.Equal(t, tc.want, got)
		}
	}
}
