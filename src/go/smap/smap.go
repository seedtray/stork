//Sorted Map is a key value container that besides Put() and Get() implements
//range queries, returning a set of elements in key order.
package smap

//Key represents both the identity for an entry in the smap and it's sort order.
type Key interface {
	Cmp(Key) int
}

//Values are just an opaque interface. Users of smap will probably use
//the same type for all values.
type Value interface{}

//Edge represents a boundary in an Interval. The Key indicates the boundary value
//and Open indicates if the interval is Open (the edge Key is considered
//outside the edge)
type Edge struct {
	Key  Key
	Open bool
}

//Inf, the Zero value edge is represents +-infinity.
var Inf Edge = Edge{}

//An interval expresses a set of keys, those within the interval edges.
//If Keys were numbers, Interval{Edge{Key(10),true}, Edge{Key(100),false}}
//represents [10,100). This is, not including 10 and including 100
//For strings, ["apple", "applf") would represent any string starting with "apple"
type Interval struct {
	From, To Edge
}

//SMap is the API of a sorted map. It comprises get, put and the scanner interface.
type SMap interface {
	SMapReader
	Put(key Key, v Value)
}

//SMapReader is a read only sorted map.
type SMapReader interface {
	Get(key Key) (v Value, found bool)
	Range(i Interval) Iterator
	Len() int
	Size() int
}

//Range() operations return an Iterator over the results.
//Next() advances the iteration one step and returns whether there
//was a next element or not.
//Key() and Value() return the key and value of the current step.
//Close() manually stops the iteration.
//
//Iterator is also a stateful protocol, and its methods semantics
//depend on the iterator state, which can be one of:
//New (initial state), Open and Closed
//
//New state:
//Key() and Value() are invalid and may panic.
//Next() transitions to Open when it retuns true,
//and to Closed when false.
//Close() transitions to Closed
//
//Open state:
//Key() and Value() are only valid within this state.
//Next() keeps the state Open if returns true,
//and transitions to Closed if false.
//
//Closed State:
//Key() and Value() are invalid and may panic.
//Next() will always return false
//Close() is a no-op, but still a valid call.
//
//
//Iterate all Example:
//	func all(s SMap, i Interval) {
//		for iter := s.Range(i); iter.Next(); {
//			do_something(iter.Key(), iter.Value())
//		}
//	}
//
//Manually stop the iteration
//	func usage2(s SMap, i Interval) {
//		for iter := s.Range(i); iter.Next(); {
//			if should_quit {
//				iter.Close()
//			} else {
//				do_something(iter.Key(), iter.Value())
//			}
//		}
//	}
//
//Wrong usage. Iterator is not closed.
//If the Iterator holds a lock or some other resource
//It won't be released.
//	func dont_do_this(s SMap, i Interval) {
//		for iter := s.Range(i); iter.Next(); {
//			if should_quit {
//				break
//			}
//			do_something(iter.Key(), iter.Value())
//		}
//	}
//
//Just use the first element
//	func one_or_nothing(s SMap, i Interval) {
//		if iter := s.Range(i); iter.Next() {
//			do_something(iter.Key(), iter.Value())
//			iter.Close()
//		}
//	}
type Iterator interface {
	Next() bool
	Key() Key
	Value() Value
	Close()
}
