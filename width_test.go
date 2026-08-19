package runewidth

import (
	"fmt"
	"testing"
)

func TestWidth(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want int
	}{
		{"negative", -1, 0},
		{"too large", 0x110000, 0},
		{"surrogate", 0xD800, 0},
		{"nul", 0, 0},
		{"control", '\x1b', 0},
		{"delete", '\x7f', 0},
		{"soft hyphen", '\u00ad', 0},
		{"combining", '\u0301', 0},
		{"enclosing", '\u0488', 0},
		{"spacing mark", '\u0903', 1},
		{"named combining spacing mark", '\u0cf3', 1},
		{"musical spacing mark", '\U0001d165', 1},
		{"zero width non-joiner", '\u200c', 0},
		{"zero width joiner", '\u200d', 0},
		{"line separator", '\u2028', 0},
		{"paragraph separator", '\u2029', 0},
		{"variation selector", '\ufe0f', 0},
		{"ascii", 'A', 1},
		{"ambiguous", '\u00a1', 1},
		{"hangul", '\u1100', 2},
		{"cjk", '界', 2},
		{"fullwidth", '\uff01', 2},
		{"emoji", '😀', 2},
		{"private use", '\ue000', 1},
		{"unassigned", '\u0378', 1},
		{"noncharacter", '\ufdd0', 0},
		{"plane noncharacter", '\U0010ffff', 0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Width(test.r); got != test.want {
				t.Fatalf("Width(%U) = %d, want %d", test.r, got, test.want)
			}
		})
	}
}

func TestWidthRange(t *testing.T) {
	for r := rune(-1); r <= 0x110000; r++ {
		w := Width(r)
		if w < 0 || w > 2 {
			t.Fatalf("Width(%U) = %d, want 0, 1, or 2", r, w)
		}
	}
}

var widthSink int

func TestWidthAllocs(t *testing.T) {
	allocs := testing.AllocsPerRun(1000, func() {
		widthSink = Width('界')
	})
	if allocs != 0 {
		t.Fatalf("Width allocated %v times, want 0", allocs)
	}
}

func ExampleWidth() {
	fmt.Println(Width('a'))
	fmt.Println(Width('\u0301'))
	fmt.Println(Width('界'))
	// Output:
	// 1
	// 0
	// 2
}
