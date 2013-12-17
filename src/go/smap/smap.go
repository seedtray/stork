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
//Next() advances the iterator one step. If it returns true, there will be a new
//Value available by calling Value()
type Iterator interface {
	Next() bool
	Key() Key
	Value() Value
}
