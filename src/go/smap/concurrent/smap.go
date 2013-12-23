//Concurrent safe implementation of smap.SMap. This is an overly simple
//implementation that composes an unsafe SMap with a read-write
//lock around its operations.
package concurrent

import (
	"github.com/losmonos/stork/src/go/smap"
	"sync"
)

//we keep the type private as it only exists to satisfy the SMap interface
type lockSMap struct {
	smap.SMap
	lock sync.RWMutex
}

func (c *lockSMap) Get(key smap.Key) (v smap.Value, found bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.SMap.Get(key)

}
func (c *lockSMap) Put(key smap.Key, value smap.Value) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.SMap.Put(key, value)
}

func (c *lockSMap) Range(i smap.Interval) smap.Iterator {
	c.lock.RLock()
	return &scanner{Iterator: c.SMap.Range(i), lock: c.lock.RLocker()}
}

//Builds a concurrent safe smap.SMap given a base SMap object.
//The only notable implementation detail of the concurrent SMap is
//that Iterators hold a read lock across the entire map, meaning
//Put() operations will block until the Iterator is either exhausted
//or manually Closed.
func New(base smap.SMap) smap.SMap {
	//also enforces lockSMap baseements SMap
	return &lockSMap{SMap: base}
}

//An iterator that holds a read lock until stopped
//Again, private as it only implements smap.Scanner
type scanner struct {
	smap.Iterator
	lock   sync.Locker
	closed bool
}

func (s *scanner) Next() bool {
	if s.closed {
		return false
	}
	if s.Iterator.Next() {
		return true
	} else {
		s.Close()
		return false
	}
}

func (s *scanner) Close() {
	if !s.closed {
		s.Iterator.Close()
		s.lock.Unlock()
		s.closed = true
	}
}

//enforce *scanner implements Iterator
var _ smap.Iterator = &scanner{}
