package redblack

import (
	"bufio"
	"fmt"
	"github.com/losmonos/stork/src/go/smap"
	"math/rand"
	"os"
	"testing"
	"testing/quick"
)

func TestEmpty(t *testing.T) {
	store := sstore{New(ssFactory)}
	checkOrderInvariant(store.SMap, t)
	if store.Len() != 0 {
		t.Error("Expected 0 length store")
	}
	if v, found := store.Get("invalid"); found {
		t.Errorf("Expected not to get any value, got invalid:%s", v)
	}
}

func checkOrderInvariant(s smap.SMap, t *testing.T) {
	m := s.(*RedBlack)
	last := str("")
	m.inOrder(func(n *Node) bool {
		current := n.entry.GetKey()
		if current.Cmp(last) < 0 {
			t.Errorf("Expected %s to be less than %s", last, n.entry.GetKey())
		}
		last = current.(str)
		return true
	})

}

func TestTwo(t *testing.T) {
	words := []string{"Hello", "Two", "Another", "AAA"}
	store := sstore{New(ssFactory)}
	for _, word := range words {
		store.Put(word, word)
	}
	checkOrderInvariant(store.SMap, t)
}

func isOrdered(s smap.SMap) bool {
	m := s.(*RedBlack)
	last := str("")
	holds := true
	m.inOrder(func(n *Node) bool {
		current := n.entry.GetKey()
		if current.Cmp(last) < 0 {
			holds = false
			return false
		}
		last = current.(str)
		return true
	})
	return holds
}

func TestOrderBig(t *testing.T) {
	store := sstore{New(ssFactory)}
	f := func(k string) bool {
		store.Put(k, "")
		return isOrdered(store.SMap)
	}
	if err := quick.Check(f, quickConfig); err != nil {
		t.Error(err)
	}
}

func gotRight(s smap.SMap) bool {
	m := s.(*RedBlack)
	holds := true
	m.inOrder(func(n *Node) bool {
		v, found := m.Get(n.entry.GetKey())
		holds = found && v == n.entry.GetValue()
		return holds
	})
	return holds
}

var quickConfig = &quick.Config{500, 0, nil, nil}

func TestGetRightValue(t *testing.T) {
	store := sstore{New(ssFactory)}
	f := func(k string, v string) bool {
		store.Put(k, v)
		return gotRight(store.SMap)
	}
	if err := quick.Check(f, quickConfig); err != nil {
		t.Error(err)
	}
}

func BenchmarkPutDistinctWords(b *testing.B) {
	b.StopTimer()
	testData.load()
	words := testData.shuffled_words
	store := sstore{New(ssFactory)}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mi := i % len(words)
		store.Put(words[mi], words[mi])
	}
}

func BenchmarkGetExistingWords(b *testing.B) {
	b.StopTimer()
	testData.load()
	words := testData.shuffled_words
	store := sstore{New(ssFactory)}
	for i := 0; i < b.N; i++ {
		mi := i % len(words)
		store.Put(words[mi], words[mi])
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mi := i % len(words)
		store.Get(words[mi])
	}
}

func BenchmarkPutDistinctInts(b *testing.B) {
	b.StopTimer()
	nums := rand.Perm(b.N)
	store := New(nnFactory)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		store.Put(number(nums[i]), i)
	}
}

func BenchmarkGetExistingInts(b *testing.B) {
	b.StopTimer()
	nums := rand.Perm(b.N)
	store := New(nnFactory)
	for i := 0; i < b.N; i++ {
		store.Put(number(nums[i]), i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		store.Get(number(nums[i]))
	}
}

type testDataCache struct {
	shuffled_words []string
	words          []string
	loaded         bool
}

var testData = testDataCache{}

func (t *testDataCache) load() {
	if t.loaded {
		return
	}
	t.shuffled_words = loadStringLines("testdata/shuffle_words")
	t.words = loadStringLines("testdata/words")
	t.loaded = true
}

func loadStringLines(filename string) []string {
	lines := make([]string, 0, 100)
	if file, err := os.Open(filename); err != nil {
		panic(fmt.Sprintf("test data not found at %s", filename))
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
	}
	return lines
}
