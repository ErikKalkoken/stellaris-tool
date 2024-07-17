package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ErikKalkoken/stellaris-tool/internal/parser"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	cases := []struct {
		in   string
		want map[string][]any
	}{
		// Regular values
		{
			"alpha=5",
			map[string][]any{"alpha": {5.0}},
		},
		{
			"alpha=5.3",
			map[string][]any{"alpha": {5.3}},
		},
		{
			"alpha=\"special text\"",
			map[string][]any{"alpha": {"special text"}},
		},
		{
			"alpha=yes",
			map[string][]any{"alpha": {true}},
		},
		{
			"alpha=no",
			map[string][]any{"alpha": {false}},
		},
		{
			"alpha=male",
			map[string][]any{"alpha": {"male"}},
		},
		// Null values
		{
			"alpha=none",
			map[string][]any{"alpha": {nil}},
		},
		{
			"alpha=not_set",
			map[string][]any{"alpha": {nil}},
		},
		// Array
		{
			"alpha={5 6}",
			map[string][]any{"alpha": {[]float64{5, 6}}},
		},
		{
			"alpha={5.1 6.2}",
			map[string][]any{"alpha": {[]float64{5.1, 6.2}}},
		},
		{
			"alpha={6.2 0}",
			map[string][]any{"alpha": {[]float64{6.2, 0}}},
		},
		{
			"alpha={0 6.2}",
			map[string][]any{"alpha": {[]float64{0, 6.2}}},
		},
		{
			"alpha={\"first\" \"second\"}",
			map[string][]any{"alpha": {[]string{"first", "second"}}},
		},
		{
			"alpha={bravo={1 2 3}}",
			map[string][]any{"alpha": {map[string][]any{"bravo": {[]float64{1, 2, 3}}}}},
		},
		{
			"alpha={{bravo=1}{bravo=2}}",
			map[string][]any{"alpha": {[]map[string][]any{{"bravo": {1.0}}, {"bravo": {2.0}}}}},
		},
		{
			"alpha={yes yes no no}",
			map[string][]any{"alpha": {[]bool{true, true, false, false}}},
		},
		// Objects
		{
			"alpha={bravo=3}",
			map[string][]any{"alpha": {map[string][]any{"bravo": {3.0}}}},
		},
		{
			"alpha={bravo=3 charlie=4}",
			map[string][]any{"alpha": {map[string][]any{"bravo": {3.0}, "charlie": {4.0}}}},
		},
		{
			"alpha=5 bravo=6 charlie=7",
			map[string][]any{"alpha": {5.0}, "bravo": {6.0}, "charlie": {7.0}},
		},
		{
			"alpha={bravo=3 charlie=7}",
			map[string][]any{"alpha": {map[string][]any{"bravo": {3.0}, "charlie": {7.0}}}},
		},
		{
			"alpha={0={bravo=1} 1={charlie=7}}",
			map[string][]any{"alpha": {
				map[string][]any{
					"0": {map[string][]any{"bravo": {1.0}}},
					"1": {map[string][]any{"charlie": {7.0}}}},
			}},
		},
		// Special cases
		{
			"alpha={}",
			map[string][]any{"alpha": {}}},
		{
			"alpha={none={}}",
			map[string][]any{"alpha": {map[string][]any{"none": {}}}},
		},
		{
			"alpha={1=\"test\"}",
			map[string][]any{"alpha": {map[string][]any{"1": {"test"}}}},
		},
		{
			"alpha={\"bravo\"=3}",
			map[string][]any{"alpha": {map[string][]any{"bravo": {3.0}}}},
		},
		{
			"alpha={1={bravo=2}}",
			map[string][]any{"alpha": {
				map[string][]any{"1": {map[string][]any{"bravo": {2.0}}}},
			}},
		},
		// Array of objects without equal sign
		{
			"alpha={{bravo 42}}",
			map[string][]any{"alpha": {[]map[string][]any{{"bravo": {42.0}}}}},
		},
		// Date as value which is no string
		{
			"alpha=2259.11.28",
			map[string][]any{"alpha": {"2259.11.28"}},
		},
		// Objects with same keys (one instance)
		{
			"alpha={bravo=3 bravo=4 bravo=9 bravo=1 bravo=2}",
			map[string][]any{"alpha": {
				map[string][]any{
					"bravo": {3.0, 4.0, 9.0, 1.0, 2.0},
				}},
			},
		},
		// Objects with same keys (multiple instances)
		{
			"alpha={bravo=3 charlie=1 bravo=4 charlie=2 bravo=9 charlie=3 bravo=1 charlie=4 bravo=2 charlie=5}",
			map[string][]any{"alpha": {
				map[string][]any{
					"bravo":   {3.0, 4.0, 9.0, 1.0, 2.0},
					"charlie": {1.0, 2.0, 3.0, 4.0, 5.0},
				}}},
		},
		// Objects with same keys (multiple instances) mixed with other k/v paris
		{
			"alpha={bravo=3 charlie=1 bravo=4 charlie=2 bravo=9 charlie=3 bravo=1 charlie=4 bravo=2 charlie=5 delta=1}",
			map[string][]any{"alpha": {
				map[string][]any{
					"bravo":   {3.0, 4.0, 9.0, 1.0, 2.0},
					"charlie": {1.0, 2.0, 3.0, 4.0, 5.0},
					"delta":   {1.0},
				}},
			}},
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

// func TestParserFull(t *testing.T) {
// 	f, err := os.Open("../.temp/gamestate")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()
// 	p := parser.NewParser(f)
// 	_, err = p.Parse()
// 	assert.NoError(t, err)
// }
