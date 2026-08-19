// Package runewidth reports terminal column widths for Unicode code points.
//
// Width is code-point based. It does not measure grapheme clusters, whose
// display width can depend on surrounding runes, terminal policy, and fonts.
package runewidth

//go:generate go run ./internal/gen -output tables.go

const (
	maxRune  = 0x10FFFF
	pageSize = 1 << 8
)

// Width reports the terminal column width of r as 0, 1, or 2.
//
// Invalid rune values, controls, format characters, combining marks, and
// Unicode noncharacters have width 0. East Asian Wide and Fullwidth runes have
// width 2. All other runes have width 1.
//
// Width performs no allocation.
func Width(r rune) int {
	u := uint32(r)
	if u > maxRune || 0xD800 <= u && u <= 0xDFFF {
		return 0
	}
	page := int(widthPage[u>>8])
	return int(widthData[page*pageSize+int(u&0xFF)])
}
