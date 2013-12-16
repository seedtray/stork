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

//Close closes the iterator. Next() will return true hereafter.
func (s *Scanner) Close() {
	s.closed = true
}

//Next advances the iterator one step and if returns true, an entry will be available upon calling Entry()
func (s *Scanner) Next() bool {
	stack := s.stack
	if s.closed || stack.Empty() {
		return false
	}
	s.node = stack.Pop()
	for node := s.node.Right; node != nil; node = node.Left {
		stack.Push(node)
	}
	return true
}

//Entry returns the current entry in the iterator.
func (s *Scanner) Value() smap.Value {
	return s.node.entry.GetValue()
}

//Scan() returns a ScannerBuilder for specifying a scan configuration in a fluent interface.
func (m *RedBlack) Scan() smap.ScannerBuilder {
	return &ScannerBuilder{m, nil, nil}
}

//ScannerBuilder is a helper for specifying different scan types
type ScannerBuilder struct {
	m           *RedBlack
	left, right *smap.Edge
}

//From specifies that the scan starts from a given key and that the range is left-closed.
func (sb *ScannerBuilder) From(k smap.Key) smap.ScannerBuilder {
	sb.left = &smap.Edge{k, true}
	return sb
}

//After specifies that the scan starts from a given key and that the range is left-open.
func (sb *ScannerBuilder) After(k smap.Key) smap.ScannerBuilder {
	sb.left = &smap.Edge{k, false}
	return sb
}

//To specifies that the scan stops at a given key and that the range is right-closed.
func (sb *ScannerBuilder) To(k smap.Key) smap.ScannerBuilder {
	sb.right = &smap.Edge{k, true}
	return sb
}

//Before specifies that the scan stops at a given key and that the range is right-open.
func (sb *ScannerBuilder) Before(k smap.Key) smap.ScannerBuilder {
	sb.right = &smap.Edge{k, false}
	return sb
}

//Start builds the scanner according to the already given configuration.
//A non configured ScannerBuilder will return a Full Scanner.
func (sb *ScannerBuilder) Start() smap.Iterator {
	m := sb.m
	if sb.left == nil && sb.right == nil {
		return m.FullScan()
	} else if sb.right == nil {
		return m.UpToScan(*sb.right)
	} else if sb.left == nil {
		return m.FromScan(*sb.left)
	} else {
		if sb.left.Key.Cmp(sb.right.Key) > 0 {
			return &Scanner{}
		} else {
			return m.RangeScan(*sb.left, *sb.right)
		}
	}
}

//upToScanner implements a scanner with a max boundary.
type upToScanner struct {
	*Scanner
	boundary *Node
	closed   bool
}

//Next advances the iterator one step and if returns true, an entry will be available upon calling Entry()
func (u *upToScanner) Next() bool {
	if u.closed {
		return false
	}
	if !u.Scanner.Next() {
		u.Close()
		return false
	} else {
		if u.Scanner.node == u.boundary {
			u.Close()
		}
		return true
	}
}

//FromScan returns an iterator that scans through the RedBlack keys in order starting at the specified edge.
func (m *RedBlack) FromScan(from smap.Edge) smap.Iterator {
	stack := m.buildLeftBoundStack(from)
	return &Scanner{stack, nil, false}
}

//FullScan returns an iterator that scans through all the RedBlack keys in order
func (m *RedBlack) FullScan() smap.Iterator {
	stack := m.makeHeightStack()
	//populate the stack with all the left wing roots
	for node := m.root; node != nil; node = node.Left {
		stack.Push(node)
	}
	return &Scanner{stack, nil, false}
}

//FromScan returns an iterator that scans through the RedBlack keys in order stopping at the specified edge.
func (m *RedBlack) UpToScan(to smap.Edge) smap.Iterator {
	boundary := m.getRightBound(to)
	scanner := m.FullScan().(*Scanner)
	return &upToScanner{scanner, boundary, false}
}

//RangeScan returns an iterator that scans through the RedBlack keys in order between the given start and end edges.
func (m *RedBlack) RangeScan(from, to smap.Edge) smap.Iterator {
	boundary := m.getRightBound(to)
	scanner := m.FromScan(from).(*Scanner)
	return &upToScanner{scanner, boundary, false}
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
			if left.Closed {
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
//The edge key does not need to be in the RedBlack
func (m *RedBlack) getRightBound(right smap.Edge) *Node {
	var last *Node = nil
	key := right.Key
	if right.Key == nil {
		return nil
	}
	for current := m.root; current != nil; {
		if cmp := current.entry.GetKey().Cmp(key); cmp == 0 {
			return current
		} else if cmp > 0 {
			current = current.Left
		} else if cmp < 0 {
			last = current
			current = current.Right
		}
	}
	return last
}
