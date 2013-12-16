//Helpers for testing llrb with different Key implementations.
//Implements string and int Keys
//Also includes a loader for test data
package llrb

import (
	"bytes"
)

type str string

//Cmp compares two string Keys
func (s str) Cmp(other Key) int {
	return bytes.Compare([]byte(string(s)), []byte(string(other.(str))))
}

//ss implements a string to string Entry
type ss struct {
	key   str
	value string
}

func (e *ss) GetKey() Key { return e.key }

func (e *ss) GetValue() Value { return e.value }

func (e *ss) SetValue(v Value) { e.value = v.(string) }

func (e *ss) Size() int { return len(e.key) + len(e.value) }

func (e *ss) Empty() bool { return false }

func ssFactory(key Key, value Value) Entry {
	return &ss{key.(str), value.(string)}
}

//enforce ss implements Entry
var _ Entry = &ss{}

//enforce ssFactory implements ssFactory
var _ EntryFactory = ssFactory

//sstore is a Memlog with str keys and str values.
//It just wraps a Smap and unwraps the Key interface value to/from strings
type sstore struct{ *Smap }

//Get returns a string value from the Smap given a string key
func (s *sstore) Get(key string) (string, bool) {
	if v, found := s.Smap.Get(str(key)); found {
		return v.(string), true
	} else {
		return "", false
	}
}

//Put saves a string value by a string key
func (s *sstore) Put(key, value string) { s.Smap.Put(str(key), value) }

//Length of the string based Smap
func (s *sstore) Len() int { return s.Smap.Len() }

//number implements a int to int Smap. Similar to sstore,
//it wraps a Smap and converts the Key values to/from strings
type number int

//nn implements a int(number) to int Entry
type nn struct {
	key   number
	value int
}

//Cmp compares to int Keys
func (n number) Cmp(other Key) int {
	return int(n) - int(other.(number))
}

func (e *nn) GetKey() Key { return e.key }

func (e *nn) GetValue() Value { return e.value }

func (e *nn) SetValue(v Value) { e.value = v.(int) }

func (e *nn) Size() int { return 16 }

func (e *nn) Empty() bool { return false }

//nnFactory implements EntryFactory for int keys and values
func nnFactory(key Key, value Value) Entry {
	return &nn{key.(number), value.(int)}
}

//enforce nnFactory implement EntryFactory
var _ EntryFactory = nnFactory
