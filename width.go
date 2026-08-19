// Package runewidth reports terminal column widths for Unicode code points.
//
// Width is code-point based. It does not measure grapheme clusters, whose
// display width can depend on surrounding runes, terminal policy, and fonts.
package runewidth

//go:generate go run ./internal/gen -output tables.go

const maxRune = 0x10FFFF

type widthRange struct {
	first rune
	last  rune
}

// Width reports the terminal column width of r as 0, 1, or 2.
//
// Invalid rune values, controls, format characters, combining marks, and
// Unicode noncharacters have width 0. East Asian Wide and Fullwidth runes have
// width 2. All other runes have width 1.
//
// Width performs no allocation.
func Width(r rune) int {
	if uint32(r) > maxRune || 0xD800 <= r && r <= 0xDFFF {
		return 0
	}
	if r < 0x300 {
		if r < 0x20 || 0x7F <= r && r <= 0x9F || r == 0xAD {
			return 0
		}
		return 1
	}
	if inTable(r, zeroWidth[:]) {
		return 0
	}
	if inTable(r, doubleWidth[:]) {
		return 2
	}
	return 1
}

func inTable(r rune, table []widthRange) bool {
	lo, hi := 0, len(table)
	for lo < hi {
		m := int(uint(lo+hi) >> 1)
		if table[m].last < r {
			lo = m + 1
		} else {
			hi = m
		}
	}
	return lo < len(table) && table[lo].first <= r
}
