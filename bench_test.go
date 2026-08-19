package runewidth

import "testing"

const benchmarkRunes = 250_000

var (
	benchmarkASCII   = repeatRunes("The quick brown fox jumps over the lazy dog. ", benchmarkRunes)
	benchmarkMixed   = repeatRunes("Go界e\u0301😀αЖ한", benchmarkRunes)
	benchmarkMillion = repeatRunes("Go界e\u0301😀αЖ한", 1_000_000)
)

func repeatRunes(s string, n int) []rune {
	pattern := []rune(s)
	runes := make([]rune, n)
	for i := range runes {
		runes[i] = pattern[i%len(pattern)]
	}
	return runes
}

func BenchmarkWidth(b *testing.B) {
	benchmarks := []struct {
		name  string
		runes []rune
	}{
		{"ASCII250K", benchmarkASCII},
		{"Mixed250K", benchmarkMixed},
		{"Mixed1M", benchmarkMillion},
	}
	for _, benchmark := range benchmarks {
		b.Run(benchmark.name, func(b *testing.B) {
			runes := benchmark.runes
			b.ReportAllocs()
			b.ReportMetric(float64(len(runes)), "runes/op")
			for b.Loop() {
				n := 0
				for _, r := range runes {
					n += Width(r)
				}
				widthSink = n
			}
		})
	}
}
