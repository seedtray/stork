//Sorted Maps interfaces
package smap

//Key represents both the identifier for an entry in the llrb and it's sort order.
type Key interface {
	Cmp(Key) int
}

//Values are just an opaque interface. Users of llrb will probably use the same type for all values.
type Value interface{}

//Edge represents a boundary in a range. The Key indicates the boundary value
//and Open indicates if the interval is Open (the edge Key is considered outside the edge)
type Edge struct {
	Key  Key
	Open bool
}

//Inf, or the Zero value edge is represents +-infinity.
var Inf Edge = Edge{}

//Interval representing the boundaries of a range scan.
type Interval struct {
	From, To Edge
}

//SMap is the api of a sorted map. It comprises get, put and the scanner interface.
type SMap interface {
	SMapReader
	Put(key Key, v Value)
}

//SMapReader is a read only SMap
type SMapReader interface {
	Get(key Key) (v Value, found bool)
	Range(i Interval) Iterator
	Len() int
	Size() int
}

//Iterator is the interface for all scan/traversal methods in llrb.
//Next() advances the iterator one step.
//Calling Value() or Key() before Next() may panic or return unmeaningful results.
//If Next() returns true, Value() will return the next element's Value and Key() will return the next element's Key.
//If Next() returns false, it means there are no more elements to iterate and Value() and Key() should not be called again.
//Calling Value() or Key() after Next() returns false may panic or return unmeaningful results.
type Iterator interface {
	Next() bool
	Key() Key
	Value() Value
}
