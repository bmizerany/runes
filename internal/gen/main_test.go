package main

import "testing"

func TestMakePages(t *testing.T) {
	width := make([]byte, 3*pageSize)
	for i := range width {
		width[i] = 1
	}
	width[17] = 2
	width[pageSize+17] = 2
	width[2*pageSize+23] = 0

	index, pages := makePages(width)
	if got, want := len(index), 3; got != want {
		t.Fatalf("len(index) = %d, want %d", got, want)
	}
	if got, want := len(pages), 2*pageSize; got != want {
		t.Fatalf("len(pages) = %d, want %d", got, want)
	}
	if index[0] != index[1] {
		t.Fatalf("equal pages have indexes %d and %d", index[0], index[1])
	}
	if index[1] == index[2] {
		t.Fatalf("unequal pages share index %d", index[1])
	}
	for i, want := range width {
		page := int(index[i/pageSize])
		if got := pages[page*pageSize+i%pageSize]; got != want {
			t.Fatalf("width[%d] = %d, want %d", i, got, want)
		}
	}
}
