//the following tests are meant to be run with the race detector.
//There are no invariant checks here. Success is determined by liveness
//and no race detection.
package concurrent

import (
	"github.com/losmonos/stork/src/go/smap"
	"github.com/losmonos/stork/src/go/smap/redblack"
	"testing"
)

type Key int

func (k Key) Cmp(other smap.Key) int {
	return int(k) - int(other.(Key))
}

type Entry struct {
	key   Key
	value bool
}

func (e *Entry) GetKey() smap.Key      { return e.key }
func (e *Entry) SetValue(v smap.Value) { e.value = v.(bool) }
func (e *Entry) GetValue() smap.Value  { return e.value }
func (e *Entry) Size() int             { return 1 }
func (e *Entry) Empty() bool           { return false }

func factory(k smap.Key, v smap.Value) redblack.Entry {
	return &Entry{k.(Key), v.(bool)}
}

//waiter is a small helper for spawning functions in goroutines and waiting
//for all of them to finish.
type waiter struct {
	howMany int
	done    chan bool
}

func newWaiter() *waiter {
	return &waiter{done: make(chan bool)}
}

//spawn f in a separate goroutine
func (w *waiter) spawn(f func()) {
	w.howMany++
	go func() {
		f()
		w.done <- true
	}()
}

//wait for all goroutines to finish.
func (w *waiter) wait() {
	for i := 0; i < w.howMany; i++ {
		<-w.done
	}
}

//mirror advance reads on a key namespace defined by the range (src,).
//waits for a key src+i to exist and then writes dst+i+offset
//two mirror_advance can work together building two incremental
//lists of set keys.
func mirror_advance(s smap.SMap, src, dst, offset, limit int) {
	for i := 0; i < limit; i++ {
		for found := false; !found; _, found = s.Get(Key(src + i)) {
		}
		s.Put(Key(dst+i+offset), true)
	}
}

//scan the entire smap over and over until there are at least _stop_ keys
//set in the smap.
func scan(s smap.SMap, stop int) {
	for found, i := 0, 0; found < stop; i++ {
		found = 0
		for it := s.Range(smap.Interval{}); it.Next(); {
			found++
		}
	}
}

//just test if the concurrent smap makes progress when mutliple
//readers and writers are active
//also make different goroutines to depend on each other side effects,
//running with test -race should not find any race conditions.
//running with test -race and using a redblack.New does find race conditions.
func TestLiveness(t *testing.T) {
	cs := New(redblack.New(factory))
	w := newWaiter()
	w.spawn(func() { mirror_advance(cs, 0, 200, 0, 20) })
	w.spawn(func() { mirror_advance(cs, 200, 0, 1, 20) })
	w.spawn(func() { scan(cs, 41) })
	cs.Put(Key(0), true)
	w.wait()
}

//benchmark a function that reads and writes to a standard redblack tree
func BenchmarkRedBlack(b *testing.B) {
	s := redblack.New(factory)
	s.Put(Key(0), true)
	mirror_advance(s, 0, 0, 1, b.N)
}

//benchmark a function that reads and writes to a concurrent redblack tree
func BenchmarkConcurrentRedBlack(b *testing.B) {
	s := New(redblack.New(factory))
	s.Put(Key(0), true)
	mirror_advance(s, 0, 0, 1, b.N)
}
