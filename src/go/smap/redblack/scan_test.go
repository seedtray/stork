package redblack

import (
	"github.com/losmonos/stork/src/go/smap"
	"sort"
	"testing"
)

func getLoadedStore(words []string) *RedBlack {
	m := New(ssFactory)
	testData.load()
	for _, word := range words {
		m.Put(str(word), word)
	}
	return m
}

var shortWordList = []string{"blueberry", "cherry", "lemon", "orange"}

func TestFullScan(t *testing.T) {
	testData.load()
	store := getLoadedStore(testData.words)
	for i, scanner := 0, store.Range(smap.Interval{}); scanner.Next(); i++ {
		scanner_word := scanner.Value().(string)
		word := testData.words[i]
		if scanner_word != word {
			t.Fatalf("Expected '%s' in FullScan, got '%s' instead'", word, scanner_word)
		}
	}
}

func TestRangeScan(t *testing.T) {
	testData.load()
	store := getLoadedStore(testData.words)
	sorted_words := sort.StringSlice(testData.words)
	start, stop := sorted_words.Search("hello"), sorted_words.Search("world")
	i := 0
	for iter := store.Range(wordInterval("hello", "world", false, false)); iter.Next(); i++ {
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

func TestEmptyFullScan(t *testing.T) {
	store := New(ssFactory)
	for scan := store.Range(smap.Interval{}); scan.Next(); {
		t.Fatalf("Empty full scan returned non empty results. got %q", scan.Value())
	}
}

func wordInterval(from, to string, left, right bool) smap.Interval {
	return smap.Interval{smap.Edge{str(from), left}, smap.Edge{str(to), right}}
}

func TestRangeOutsideRight(t *testing.T) {
	store := getLoadedStore(shortWordList)
	if scan := store.Range(wordInterval("pear", "tangerine", false, false)); scan.Next() {
		t.Fatalf("Empty range scan returned non empty results. got %q", scan.Value())
	}
}

func TestRangeOutsideLeft(t *testing.T) {
	store := getLoadedStore(shortWordList)
	if scan := store.Range(wordInterval("apple", "banana", false, false)); scan.Next() {
		t.Fatalf("Empty range scan returned non empty results. got %q", scan.Value())
	}
}

func TestRangeLesserToValue(t *testing.T) {
	store := getLoadedStore(shortWordList)
	scan := store.Range(wordInterval("apple", "cherry", false, false))
	scan.Next()
	if got := scan.Value(); got != "blueberry" {
		t.Fatalf("Expected 'blueberry' in scan, got %q", got)
	}
	scan.Next()
	if got := scan.Value(); got != "cherry" {
		t.Fatalf("Expected 'cherry' in scan, got %q", got)
	}
	if scan.Next() != false {
		t.Fatalf("Expected end of iteration, got %q", scan.Value())
	}
}

func TestRangeValueToGreater(t *testing.T) {
	store := getLoadedStore(shortWordList)
	scan := store.Range(wordInterval("lemon", "pear", false, false))
	scan.Next()
	if got := scan.Value(); got != "lemon" {
		t.Fatalf("Expected 'lemon' in scan, got %q", got)
	}
	scan.Next()
	if got := scan.Value(); got != "orange" {
		t.Fatalf("Expected 'orange' in scan, got %q", got)
	}
	if scan.Next() != false {
		t.Fatalf("Expected end of iteration, got %q", scan.Value())
	}

}

func TestRangeLesserToGreater(t *testing.T) {
	store := getLoadedStore(shortWordList)
	scan := store.Range(wordInterval("apple", "pear", false, false))
	i := 0
	for ; scan.Next(); i++ {
		if expect, got := shortWordList[i], scan.Value(); got != expect {
			t.Fatalf("Expected '%s' in scan, got %q", expect, got)
		}
	}
	if expected := len(shortWordList); i != expected {
		t.Fatalf("Expected %d results, got %d", expected, i)
	}
}

func TestEmptyOpenIntervals(t *testing.T) {
	store := getLoadedStore(shortWordList)
	msg := "Empty Open Interval returned non empty results, got %q"
	if scan := store.Range(wordInterval("cherry", "cherry", true, true)); scan.Next() {
		t.Fatalf(msg, scan.Value())
	}
	if scan := store.Range(wordInterval("apple", "apple", true, true)); scan.Next() {
		t.Fatalf(msg, scan.Value())
	}
	if scan := store.Range(wordInterval("pear", "pear", true, true)); scan.Next() {
		t.Fatalf(msg, scan.Value())
	}
}

func TestOpenInterval(t *testing.T) {
	store := getLoadedStore(shortWordList)
	scan := store.Range(wordInterval("cherry", "orange", true, true))
	if expect := "lemon"; !scan.Next() {
		t.Fatalf("Expected '%s', got nothing instead", expect)
	} else if got := scan.Value(); got != expect {
		t.Fatalf("Expected '%s', got '%s' instead", expect, got)
	}
	if scan.Next() {
		t.Fatalf("Expected end of iteration, got %q instead", scan.Value())
	}
}

func TestInvertedInterval(t *testing.T) {
	store := getLoadedStore(shortWordList)
	msg := "Inverted Interval returned non empty results, got %q"
	if scan := store.Range(wordInterval("orange", "cherry", true, true)); scan.Next() {
		t.Fatalf(msg, scan.Value())
	}
	if scan := store.Range(wordInterval("orange", "cherry", true, false)); scan.Next() {
		t.Fatalf(msg, scan.Value())
	}
	if scan := store.Range(wordInterval("orange", "cherry", false, true)); scan.Next() {
		t.Fatalf(msg, scan.Value())
	}
	if scan := store.Range(wordInterval("orange", "cherry", false, false)); scan.Next() {
		t.Fatalf(msg, scan.Value())
	}
}
