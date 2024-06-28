package parser_test

import (
	"example/stellaris-tool/parser"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	cases := []struct {
		in   string
		want map[string]any
	}{
		// simple values
		{"alpha=5", map[string]any{"alpha": 5}},
		{"alpha=5.3", map[string]any{"alpha": 5.3}},
		{"alpha=\"special text\"", map[string]any{"alpha": "special text"}},
		{"alpha=yes", map[string]any{"alpha": true}},
		{"alpha=no", map[string]any{"alpha": false}},
		{"alpha=none", map[string]any{"alpha": parser.None}},
		{"alpha=not_set", map[string]any{"alpha": parser.NotSet}},
		{"alpha=male", map[string]any{"alpha": "male"}},
		// ID object
		{"1=\"test\"", map[string]any{"1": "test"}},
		// Array
		{"alpha={5 6}", map[string]any{"alpha": []int{5, 6}}},
		{"alpha={5.1 6.2}", map[string]any{"alpha": []float64{5.1, 6.2}}},
		{"alpha={0 6.2}", map[string]any{"alpha": []float64{0, 6.2}}},
		{"alpha={1 2}", map[string]any{"alpha": []int{1, 2}}},
		{"alpha={\"first\" \"second\"}", map[string]any{"alpha": []string{"first", "second"}}},
		{"alpha={bravo={1 2 3}}", map[string]any{"alpha": map[string]any{"bravo": []int{1, 2, 3}}}},
		{
			"alpha={{bravo=1}{bravo=2}}",
			map[string]any{"alpha": []map[string]any{{"bravo": 1}, {"bravo": 2}}},
		},
		{"alpha={yes yes no no}", map[string]any{"alpha": []bool{true, true, false, false}}},
		// Objects
		{"alpha={bravo=3}", map[string]any{"alpha": map[string]any{"bravo": 3}}},
		{"alpha=5 bravo=6 charlie=7", map[string]any{"alpha": 5, "bravo": 6, "charlie": 7}},
		{
			"alpha={bravo=3 charlie=7}",
			map[string]any{"alpha": map[string]any{"bravo": 3, "charlie": 7}},
		},
		{
			"alpha={0={bravo=1} 1={charlie=7}}",
			map[string]any{"alpha": map[string]any{
				"0": map[string]any{"bravo": 1},
				"1": map[string]any{"charlie": 7}},
			},
		},
		// Empty object
		{"alpha={}", map[string]any{"alpha": struct{}{}}},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.in), func(t *testing.T) {
			r := strings.NewReader(tc.in)
			p := parser.NewParser(r)
			got, err := p.Parse()
			if assert.NoError(t, err) {
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestParserFull(t *testing.T) {
	f, err := os.Open("testdata/example")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	p := parser.NewParser(f)
	_, err = p.Parse()
	assert.NoError(t, err)
}
