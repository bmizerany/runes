package runes

import (
	"testing"

	mattn "github.com/mattn/go-runewidth"
)

var mattnWidth = &mattn.Condition{StrictEmojiNeutral: true}

func FuzzWidthAgainstMattn(f *testing.F) {
	for _, r := range []rune{
		-1, 0, '\x1b', '\x7f', '\u00ad', '\u0301', '\u0488',
		'\u0903', '\u0cf3', '\u200c', '\u200d', '\u2028', '\u2029',
		'\ufe0f', '\U0001d165', 'A', '\u00a1',
		'\u1100', '界', '\uff01', '😀', '\ue000', '\u0378',
		'\ufdd0', rune(0xD800), maxRune, maxRune + 1,
	} {
		f.Add(int32(r))
	}

	f.Fuzz(func(t *testing.T, value int32) {
		r := rune(value)
		got, want := Width(r), mattnWidth.RuneWidth(r)
		if got != want && classifyMattnDifference(r, got, want) == mattnUnknown {
			t.Fatalf("Width(%U) = %d; mattn RuneWidth = %d", r, got, want)
		}
	})
}

func TestWidthAgainstMattn(t *testing.T) {
	var counts [3]int
	for r := rune(-1); r <= maxRune+1; r++ {
		got, want := Width(r), mattnWidth.RuneWidth(r)
		if got == want {
			counts[mattnSame]++
			continue
		}
		kind := classifyMattnDifference(r, got, want)
		if kind == mattnUnknown {
			t.Fatalf("Width(%U) = %d; mattn RuneWidth = %d", r, got, want)
		}
		counts[kind]++
	}
	want := [3]int{1_113_894, 211, 9}
	if counts != want {
		t.Fatalf("comparison counts = %v, want %v", counts, want)
	}
	t.Logf("same=%d policy=%d mattn-name-bug=%d", counts[0], counts[1], counts[2])
}

type mattnDifference uint8

const (
	mattnSame mattnDifference = iota
	mattnPolicy
	mattnNameBug
	mattnUnknown
)

func classifyMattnDifference(r rune, got, want int) mattnDifference {
	switch {
	case got == 0 && want == 1 && (mattnFormatPolicyGap(r) || mattnNoncharacterPolicyGap(r)):
		return mattnPolicy
	case got == 1 && want == 0 && mattnCombiningNameBug(r):
		return mattnNameBug
	}
	return mattnUnknown
}

// mattnFormatPolicyGap is the complete set of Cf runes for which v0.0.28
// applies its default width 1 rather than this package's all-Cf-zero policy.
func mattnFormatPolicyGap(r rune) bool {
	return 0x0600 <= r && r <= 0x0605 ||
		r == 0x061C || r == 0x06DD ||
		0x0890 <= r && r <= 0x0891 || r == 0x08E2 ||
		0x2060 <= r && r <= 0x2064 ||
		0x2066 <= r && r <= 0x2069 ||
		r == 0x110BD || r == 0x110CD ||
		0x13430 <= r && r <= 0x1343F ||
		0x1BCA0 <= r && r <= 0x1BCA3 ||
		0x1D173 <= r && r <= 0x1D17A ||
		r == 0xE0001 ||
		0xE0020 <= r && r <= 0xE007F
}

// mattnNoncharacterPolicyGap is the complete set of noncharacters for which
// v0.0.28 applies its default width 1. It already treats U+FFFE and U+FFFF as
// nonprinting, so the terminal pair starts at plane 1 here.
func mattnNoncharacterPolicyGap(r rune) bool {
	return 0xFDD0 <= r && r <= 0xFDEF ||
		0x1FFFE <= r && r <= maxRune && r&0xFFFE == 0xFFFE
}

// mattnCombiningNameBug is the complete set of Unicode 17 Mc spacing marks
// that v0.0.28 gives width 0. Its generator says Mc occupies one cell, but
// also selects any EastAsianWidth.txt line whose character name contains
// "COMBINING", which admits these nine Mc runes.
func mattnCombiningNameBug(r rune) bool {
	return r == 0x0CF3 ||
		0x1D165 <= r && r <= 0x1D166 ||
		0x1D16D <= r && r <= 0x1D172
}
