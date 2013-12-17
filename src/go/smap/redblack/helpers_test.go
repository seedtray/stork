//Helpers for testing SMaps with different Key implementations.
//Implements string and int Keys
//Also includes a loader for test data
//Also implements an in-order node traversal for invariant checking.
package redblack

import (
	"bytes"
	"github.com/losmonos/stork/src/go/smap"
)

type str string

//Cmp compares two string Keys
func (s str) Cmp(other smap.Key) int {
	return bytes.Compare([]byte(string(s)), []byte(string(other.(str))))
}

//ss implements a string to string Entry
type ss struct {
	key   str
	value string
}

func (e *ss) GetKey() smap.Key { return e.key }

func (e *ss) GetValue() smap.Value { return e.value }

func (e *ss) SetValue(v smap.Value) { e.value = v.(string) }

func (e *ss) Size() int { return len(e.key) + len(e.value) }

func (e *ss) Empty() bool { return false }

func ssFactory(key smap.Key, value smap.Value) Entry {
	return &ss{key.(str), value.(string)}
}

//enforce ss implements Entry
var _ Entry = &ss{}

//enforce ssFactory implements ssFactory
var _ EntryFactory = ssFactory

//sstore is a Memlog with str keys and str values.
//It just composes a SMap and unwraps the Key interface value to/from strings
type sstore struct{ smap.SMap }

//Get returns a string value from the SMap given a string key
func (s *sstore) Get(key string) (string, bool) {
	if v, found := s.SMap.Get(str(key)); found {
		return v.(string), true
	} else {
		return "", false
	}
}

//Put saves a string value by a string key
func (s *sstore) Put(key, value string) { s.SMap.Put(str(key), value) }

//Length of the string based SMap
func (s *sstore) Len() int { return s.SMap.Len() }

//number implements a int to int SMap. Similar to sstore,
//it wraps a SMap and converts the Key values to/from strings
type number int

//nn implements a int(number) to int Entry
type nn struct {
	key   number
	value int
}

//Cmp compares to int Keys
func (n number) Cmp(other smap.Key) int {
	return int(n) - int(other.(number))
}

func (e *nn) GetKey() smap.Key { return e.key }

func (e *nn) GetValue() smap.Value { return e.value }

func (e *nn) SetValue(v smap.Value) { e.value = v.(int) }

func (e *nn) Size() int { return 16 }

func (e *nn) Empty() bool { return false }

//nnFactory implements EntryFactory for int keys and values
func nnFactory(key smap.Key, value smap.Value) Entry {
	return &nn{key.(number), value.(int)}
}

//enforce nnFactory implement EntryFactory
var _ EntryFactory = nnFactory

//Visitor is a function that receives a Node and returns a boolean,
//true means the tree traversal should continue
type visitor func(n *Node) bool

//An in-order traversal of a Tree Node.
func (n *Node) inOrder(visit visitor) {
	if n == nil {
		return
	}
	n.Left.inOrder(visit)
	visit(n)
	n.Right.inOrder(visit)
}

//An in-order traversal of a RedBlack
func (m *RedBlack) inOrder(f visitor) {
	m.root.inOrder(f)
}
