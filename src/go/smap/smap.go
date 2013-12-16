//Sorted Maps interfaces
package smap

//Key represents both the identifier for an entry in the llrb and it's sort order.
type Key interface {
	Cmp(Key) int
}

//Values are just an opaque interface. Users of llrb will probably use the same type for all values.
type Value interface{}

//Edge represents a boundary in a range. The Key indicates the boundary value
//and Closed indicates if that key should be included in the range or not.
type Edge struct {
	Key    Key
	Closed bool
}

//SMap is the api of a sorted map. It comprises get, put and the scanner interface.
type SMap interface {
	SMapReader
	Put(key Key, v Value)
}

//SMapReader is a read only SMap
type SMapReader interface {
	Get(key Key) (v Value, found bool)
	Len() int
	Size() int
	Scanner
}

//Scanners allow range queries over an Smap, returning an Iterator of values.
type Scanner interface {
	Scan() ScannerBuilder
	FromScan(from Edge) Iterator
	UpToScan(upto Edge) Iterator
	RangeScan(from Edge, to Edge) Iterator
	FullScan() Iterator
}

//Iterator is the interface for all scan/traversal methods in llrb.
//Next() advances the iterator one step. If it returns true, there will be a new
//Entry available by calling Entry()
//Close() will stop the iterator and cause Next() to return always false from then on.
type Iterator interface {
	Next() bool
	Value() Value
	Close()
}

//ScannerBuilder is a helper for building range scans, with a fluent interface.
type ScannerBuilder interface {
	From(key Key) ScannerBuilder
	To(key Key) ScannerBuilder
	After(key Key) ScannerBuilder
	Before(key Key) ScannerBuilder
	Start() Iterator
}
