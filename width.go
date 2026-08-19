// Package runes reports terminal column widths for Unicode code points.
//
// Width is code-point based. It does not measure grapheme clusters, whose
// display width can depend on surrounding runes, terminal policy, and fonts.
//
// This package is deliberately narrower than Yasuhiro Matsumoto's excellent
// github.com/mattn/go-runewidth, which provides string width, wrapping,
// truncation, padding, and configurable conditions. Its implementation and
// Unicode tables were invaluable references in defining and auditing Width.
//
// Width uses a generated lookup table and performs no allocation. On an Apple
// M4 Pro, the package benchmark processes one million mixed code points in
// about 0.77 ms, about 39% less time than github.com/mattn/go-runewidth v0.0.28
// under non-East-Asian, strict-emoji-neutral settings. Results vary by system.
package runes

//go:generate go run ./internal/gen -output tables.go

const (
	maxRune        = 0x10FFFF
	pageSize       = 1 << 8
	packedPageSize = pageSize / 4
)

// Width reports the terminal column width of r as 0, 1, or 2.
//
// Invalid rune values, controls, format characters, line and paragraph
// separators, nonspacing and enclosing marks, and Unicode noncharacters have
// width 0. East Asian Wide and Fullwidth runes have width 2. All other runes,
// including spacing marks, have width 1.
//
// Width performs no allocation.
func Width(r rune) int {
	u := uint32(r)
	if u > maxRune || 0xD800 <= u && u <= 0xDFFF {
		return 0
	}
	page := int(widthPage[u>>8])
	i := page*packedPageSize + int(u&0xFF)/4
	shift := (u & 3) * 2
	return int(widthData[i] >> shift & 3)
}
