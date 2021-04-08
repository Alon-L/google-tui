package search

import "testing"

func TestSearcher(t *testing.T) {
	searcher := NewSearcher()

	err := searcher.Search("hello")
	if err != nil {
		t.Error(err)
	}
}
