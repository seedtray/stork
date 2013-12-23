package redblack

import (
	"github.com/losmonos/stork/src/go/smap"
	"math"
)

//fixedNodeStack is a simple fixed size stack of *Node, used in iterative tree traversals.
type fixedNodeStack struct {
	nodes []*Node
	next  int
}

//Push pushes one node on top of the stack
func (stack *fixedNodeStack) Push(node *Node) {
	stack.nodes[stack.next] = node
	stack.next++
}

//Pop removes the topmost node and returns it
func (stack *fixedNodeStack) Pop() *Node {
	stack.next--
	node := stack.nodes[stack.next]
	stack.nodes[stack.next] = nil
	return node
}

//Peek returns the topmost node without removing it.
func (stack *fixedNodeStack) Peek() *Node {
	return stack.nodes[stack.next-1]

}

//Empty tells whether there are no more nodes in the stack.
func (stack *fixedNodeStack) Empty() bool {
	return stack.next == 0
}

//maxHeight returns the maximum possible tree height as long as no new nodes are added to it.
func (m *RedBlack) maxHeight() int {
	return 2 * int(math.Ceil(math.Log2(float64(m.Len()))))
}

//Scanner iterates over all the nodes initially pushed onto its stack along with their right subtrees
//The order of iteration is: root (the node in the stack), then in-order iteration of the node's right subtree
//If the stack is populated with the leftmost walk of a binary tree,
//then this implements an in-order full scan iterator
//It is meant to fullfill both a full scan and a scan starting from a given value.
type Scanner struct {
	stack  *fixedNodeStack
	node   *Node
	closed bool
}

type EmptyScanner struct{}

func (s EmptyScanner) Next() bool {
	return false
}
func (s EmptyScanner) Value() smap.Value {
	panic("Empty Scanner")
}
func (s EmptyScanner) Key() smap.Key {
	panic("Empty Scanner")
}

func (s EmptyScanner) Close() {}

//Next advances the iterator one step and if returns true, an entry will be available upon calling Entry()
func (s *Scanner) Next() bool {
	if s.closed {
		return false
	}
	stack := s.stack
	if stack.Empty() {
		s.Close()
		return false
	}
	s.node = stack.Pop()
	for node := s.node.Right; node != nil; node = node.Left {
		stack.Push(node)
	}
	return true
}

//Value returns the current value in the iterator.
func (s *Scanner) Value() smap.Value {
	return s.node.entry.GetValue()
}

//Key returns the current key in the iterator.
func (s *Scanner) Key() smap.Key {
	return s.node.entry.GetKey()
}

func (s *Scanner) Close() {
	s.stack = nil
	s.node = nil
	s.closed = true
}

//Scan() returns an Iterator that iterates over the tree elements in order within the given interval
func (m *RedBlack) Range(i smap.Interval) smap.Iterator {
	if i.From == smap.Inf && i.To == smap.Inf {
		return m.fullScan()
	} else if i.To == smap.Inf {
		return m.upToScan(i.To)
	} else if i.From == smap.Inf {
		return m.fromScan(i.From)
	} else {
		if i.From.Key.Cmp(i.To.Key) > 0 {
			return EmptyScanner{}
		} else {
			return m.rangeScan(i.From, i.To)
		}
	}
}

//upToScanner implements a scanner with a max boundary.
type upToScanner struct {
	*Scanner
	boundary     *Node
	hit_boundary bool
	closed       bool
}

//Next advances the iterator one step and if returns true, an entry will be available upon calling Entry()
func (u *upToScanner) Next() bool {
	if u.closed {
		return false
	}
	if u.hit_boundary || !u.Scanner.Next() {
		u.Close()
		return false
	} else {
		if u.Scanner.node == u.boundary {
			u.hit_boundary = true
		}
		return true
	}
}

//FromScan returns an iterator that scans through the RedBlack keys in order starting at the specified edge.
func (m *RedBlack) fromScan(from smap.Edge) smap.Iterator {
	stack := m.buildLeftBoundStack(from)
	return &Scanner{stack: stack}
}

//FullScan returns an iterator that scans through all the RedBlack keys in order
func (m *RedBlack) fullScan() smap.Iterator {
	stack := m.makeHeightStack()
	//populate the stack with all the left wing roots
	for node := m.root; node != nil; node = node.Left {
		stack.Push(node)
	}
	return &Scanner{stack: stack}
}

//FromScan returns an iterator that scans through the RedBlack keys in order stopping at the specified edge.
func (m *RedBlack) upToScan(to smap.Edge) smap.Iterator {
	boundary := m.getRightBound(to)
	if boundary == nil {
		return EmptyScanner{}
	} else {
		scanner := m.fullScan().(*Scanner)
		return &upToScanner{Scanner: scanner, boundary: boundary}
	}
}

//RangeScan returns an iterator that scans through the RedBlack keys in order between the given start and end edges.
func (m *RedBlack) rangeScan(from, to smap.Edge) smap.Iterator {
	boundary := m.getRightBound(to)
	if boundary == nil {
		return EmptyScanner{}
	}
	scanner := m.fromScan(from).(*Scanner)
	//abort if scanner stack is empty, further checks need a non-empty stack.
	if scanner.stack.Empty() {
		return EmptyScanner{}
	}
	rightKey := boundary.entry.GetKey()
	leftKey := scanner.stack.Peek().entry.GetKey()
	if leftKey.Cmp(rightKey) > 0 {
		return EmptyScanner{}
	}
	return &upToScanner{Scanner: scanner, boundary: boundary}
}

//allocate a stack suitable for traversing the tree.
func (m *RedBlack) makeHeightStack() *fixedNodeStack {
	return &fixedNodeStack{make([]*Node, m.maxHeight()), 0}
}

//build a stack with the left-most list of disjoint sub-trees that together will satisfy:
//a Scanner will traverse the nodes in order.
//the Scanner will visit all (and only) the nodes which keys are >= left edge
func (m *RedBlack) buildLeftBoundStack(left smap.Edge) *fixedNodeStack {
	stack := m.makeHeightStack()
	key := left.Key
	for current := m.root; current != nil; {
		if cmp := current.entry.GetKey().Cmp(key); cmp == 0 {
			if left.Open {
				for current = current.Right; current != nil; current = current.Left {
					stack.Push(current)
				}
			} else {
				stack.Push(current)
			}
			break
		} else if cmp > 0 {
			stack.Push(current)
			current = current.Left
		} else if cmp < 0 {
			current = current.Right
		}
	}
	return stack
}

//find the right boundary in the RedBlack given a right edge.
//The returning node will have the greater key
//such that node.Key <= right.Key (< for open edge).
//If returns nil, it means the right edge is less than all the keys in the tree
//The edge key does not need to be in the RedBlack
func (m *RedBlack) getRightBound(right smap.Edge) *Node {
	var last *Node = nil
	key := right.Key
	for current := m.root; current != nil; {
		if cmp := current.entry.GetKey().Cmp(key); cmp == 0 {
			if right.Open {
				for current = current.Left; current != nil; current = current.Right {
					last = current
				}
				return last
			} else {
				return current
			}
		} else if cmp > 0 {
			current = current.Left
		} else if cmp < 0 {
			last = current
			current = current.Right
		}
	}
	return last
}
