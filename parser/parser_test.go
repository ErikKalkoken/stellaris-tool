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
		{"1=\"test\"", map[string]any{"1": "test"}},
		{"alpha={5 6}", map[string]any{"alpha": []int{5, 6}}},
		{"alpha={\"first\" \"second\"}", map[string]any{"alpha": []string{"first", "second"}}},
		{"alpha={bravo=3}", map[string]any{"alpha": map[string]any{"bravo": 3}}},
		{"alpha={bravo={1 2 3}}", map[string]any{"alpha": map[string]any{"bravo": []int{1, 2, 3}}}},
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

// func TestParserFull(t *testing.T) {
// 	f, err := os.Open("testdata/meta")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()
// 	p := parser.NewParser(f)
// 	got, err := p.Parse()
// 	if assert.NoError(t, err) {
// 		fmt.Println(got)
// 	}
// }
