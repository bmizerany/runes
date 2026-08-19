# runes

Package `runes` reports the terminal column width of a Unicode code
point. It has one operation:

```go
func Width(rune) int
```

`Width` returns 0, 1, or 2 and performs no allocation.

## Install

```sh
go get blake.io/runes
```

## Use

```go
package main

import (
	"fmt"

	"blake.io/runes"
)

func main() {
	fmt.Println(runes.Width('a'))      // 1
	fmt.Println(runes.Width('\u0301')) // 0
	fmt.Println(runes.Width('界'))      // 2
}
```

Invalid runes, controls, format characters, line and paragraph separators,
nonspacing and enclosing marks, and Unicode noncharacters have width 0.
East Asian Wide and Fullwidth runes have width 2. Everything else,
including spacing marks, has width 1.

The operation is deliberately code-point based. It does not measure grapheme
clusters, whose display width can depend on neighboring runes, terminal
policy, and fonts.

## Credit

Yasuhiro Matsumoto's
[`go-runewidth`](https://github.com/mattn/go-runewidth) is an excellent,
full-featured library for terminal-width work. It handles string width,
wrapping, truncation, padding, and configurable conditions—far more than this
package's one operation. Its implementation and Unicode tables were invaluable
references while defining and auditing this package. If you need that broader
API, use it. Thank you, mattn.

## Performance

`Width` uses a generated lookup table and performs no allocation.
`BenchmarkWidthAgainstMattn` processes one million mixed code points. Ten
300 ms samples on darwin/arm64, Apple M4 Pro, produced:

| Implementation | Time/op | B/op | Allocs/op |
| --- | ---: | ---: | ---: |
| `blake.io/runes` | 768.5 µs | 0 | 0 |
| `github.com/mattn/go-runewidth` v0.0.28 | 1260.4 µs | 0 | 0 |

`blake.io/runes` used 39.03% less time in this benchmark. The mattn
condition is non-East-Asian with `StrictEmojiNeutral` enabled. This compares
lookup cost, not identical width policy; the packages intentionally differ on
some code points. Results vary by system.

Run the benchmark with:

```sh
go test -run '^$' -bench '^BenchmarkWidthAgainstMattn$' -benchmem -count=10 -benchtime=300ms
```

The generated table is built from pinned Unicode 17 data:

```sh
go generate ./...
go test ./...
```

## License

BSD 3-Clause, the same license used by Go. See `LICENSE`.
