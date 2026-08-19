package runewidth

import (
	"testing"

	mattn "github.com/mattn/go-runewidth"
)

const benchmarkRunes = 250_000

var (
	benchmarkASCII   = repeatRunes("The quick brown fox jumps over the lazy dog. ", benchmarkRunes)
	benchmarkMixed   = repeatRunes("Go界e\u0301😀αЖ한", benchmarkRunes)
	benchmarkMillion = repeatRunes("Go界e\u0301😀αЖ한", 1_000_000)
	benchmarkMattn   = &mattn.Condition{StrictEmojiNeutral: true}
)

func repeatRunes(s string, n int) []rune {
	pattern := []rune(s)
	runes := make([]rune, n)
	for i := range runes {
		runes[i] = pattern[i%len(pattern)]
	}
	return runes
}

func BenchmarkWidthAgainstMattn(b *testing.B) {
	benchmarkMattn.RuneWidth('界') // Build mattn's lazy table before timing.

	b.Run("runewidth/Mixed1M", func(b *testing.B) {
		b.ReportAllocs()
		b.ReportMetric(float64(len(benchmarkMillion)), "runes/op")
		for b.Loop() {
			n := 0
			for _, r := range benchmarkMillion {
				n += Width(r)
			}
			widthSink = n
		}
	})
	b.Run("mattn/Mixed1M", func(b *testing.B) {
		b.ReportAllocs()
		b.ReportMetric(float64(len(benchmarkMillion)), "runes/op")
		for b.Loop() {
			n := 0
			for _, r := range benchmarkMillion {
				n += benchmarkMattn.RuneWidth(r)
			}
			widthSink = n
		}
	})
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
