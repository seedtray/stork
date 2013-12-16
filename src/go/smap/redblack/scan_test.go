package redblack

import (
	"sort"
	"testing"
)

func getLoadedStore() *RedBlack {
	m := New(ssFactory)
	testData.load()
	for _, word := range testData.words {
		m.Put(str(word), word)

	}
	return m
}

func TestFullScan(t *testing.T) {
	store := getLoadedStore()
	testData.load()
	for i, scanner := 0, store.Scan().Start(); scanner.Next(); i++ {
		scanner_word := scanner.Value().(string)
		word := testData.words[i]
		if scanner_word != word {
			t.Fatalf("Expected '%s' in FullScan, got '%s' instead'", word, scanner_word)
		}
	}

}

func TestRangeScan(t *testing.T) {
	store := getLoadedStore()
	testData.load()
	sorted_words := sort.StringSlice(testData.words)
	start, stop := sorted_words.Search("hello"), sorted_words.Search("world")
	i := 0
	for iter := store.Scan().From(str("hello")).To(str("world")).Start(); iter.Next(); i++ {
		iter_word := iter.Value().(string)
		word := testData.words[start+i]
		if iter_word != word {
			t.Fatalf("Expected '%s' in Range Scanner, got '%s' instead'", word, iter_word)
		}
	}
	if i+start != stop+1 {
		t.Fatalf("Range scanner didn't stop at [%d]='%s', stopped at [%d]='%s'",
			stop, testData.words[stop], i+start, testData.words[i+start])
	}
}
