package llrb

//Key represents both the identifier for an entry in the Smap and it's sort order.
type Key interface {
	Cmp(Key) int
}

//Values are just an opaque interface. Users of Smap will probably use the same type for all values.
type Value interface{}

type Entry interface {
	GetKey() Key
	SetValue(v Value)
	GetValue() Value
	Size() int
	Empty() bool
}

//EntryFactory builds an entry given a Key and a Value
type EntryFactory func(Key, Value) Entry

//A Node is the main element in the llrb structure.
//It holds an entry, links to its childs and the red/black color
type Node struct {
	entry       Entry
	Left, Right *Node
	Color       bool
}

//XXX constants should be upper case, but what about unexported constants?
const (
	red   = false
	black = true
)

//NewSmap creates a new Smap
func New(factory EntryFactory) *Smap {
	return &Smap{nil, factory, 0, 0}
}

//Smap implements a sorted Map
type Smap struct {
	root    *Node
	factory EntryFactory
	length  int
	bytes   int
}

//isRed returns whether the node is red (or black). Nil nodes (empty leafs) are black
func (n *Node) isRed() bool {
	return n != nil && n.Color == red
}

//colorFlip inverts the colors on a node and it's childs
func (n *Node) colorFlip() {
	n.Color = !n.Color
	n.Left.Color = !n.Left.Color
	n.Right.Color = !n.Right.Color
}

//rotateLeft does an anti-clockwise node rotation
func (n *Node) rotateLeft() *Node {
	x := n.Right
	n.Right = x.Left
	x.Left = n
	x.Color = n.Color
	n.Color = red
	return x

}

//rotateRight does a clockwise node rotation
func (n *Node) rotateRight() *Node {
	x := n.Left
	n.Left = x.Right
	x.Right = n
	x.Color = n.Color
	n.Color = red
	return x
}

//Len returns the amount of non empty nodes in a Smap
func (m *Smap) Len() int {
	return m.length
}

//Get searches for a given key and returns it's associated value
//and a boolean indicating if it was found
func (m *Smap) Get(key Key) (v Value, found bool) {
	for current := m.root; current != nil; {
		if cmp := current.entry.GetKey().Cmp(key); cmp == 0 {
			return current.entry.GetValue(), true
		} else if cmp > 0 {
			current = current.Left
		} else if cmp < 0 {
			current = current.Right
		}
	}
	return nil, false
}

//Put inserts a value identified by a key. if the key already existed,
//Entry.SetValue(v) will be called on the already occupied slot.
//It's up to Entry implementation to define whether the value is replaced
//or some other action is taken. It may be possible to build a multi-map
//by having Entry store the values in a collection.
//Note that this may lead to a more complicate semantic of delete() which
//is not yet implemented.
func (m *Smap) Put(key Key, value Value) {
	m.root = m.insert(m.root, key, value)
	m.root.Color = black
}

//insert does the left-leaning red black tree rotations and color flips.
//It returns the new tree root.
func (m *Smap) insert(node *Node, key Key, value Value) *Node {
	if node == nil {
		entry := m.factory(key, value)
		node := &Node{entry, nil, nil, red}
		m.length++
		m.bytes += entry.Size()
		return node
	}
	if node.Left.isRed() && node.Right.isRed() {
		node.colorFlip()
	}
	if cmp := node.entry.GetKey().Cmp(key); cmp == 0 {
		m.bytes -= node.entry.Size()
		node.entry.SetValue(value)
		m.bytes += node.entry.Size()
	} else if cmp < 0 {
		node.Right = m.insert(node.Right, key, value)
	} else if cmp > 0 {
		node.Left = m.insert(node.Left, key, value)
	}
	if node.Right.isRed() && !node.Left.isRed() {
		node = node.rotateLeft()
	}
	if node.Left.isRed() && node.Left.Left.isRed() {
		node = node.rotateRight()
	}
	return node
}

//Visitor is a function that receives a Node and returns a boolean,
//true means the tree traversal should continue
//this may go to helpers_test if it's not useful for the exported implementation.
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

//An in-order traversal of a Smap
func (m *Smap) inOrder(f visitor) {
	m.root.inOrder(f)
}
