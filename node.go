package otree

import (
	"fmt"
	"strings"
)

// Node represents internal and external nodes. Internal nodes have children,
// external nodes don't. A node can hold data.
type Node struct {
	data     interface{} // data
	parent   *Node       // parent node
	siblings []*Node     // sibling nodes
}

// WalkFunc is a function type that can be performed on all nodes in a
// (sub)tree starting at root.
type WalkFunc func(root *Node, data interface{})

type dummyType struct{}

var dummy = dummyType{}

// Ancestors returns the node's ancestors. The first one is it parent, the next
// one its grandparent and so on.
func (nd *Node) Ancestors() []*Node {
	ancestors := []*Node{}

	for p := nd.parent; p != nil; p = p.parent {
		ancestors = append(ancestors, p)
	}
	return ancestors
}

// Degree returns the node's degree, i.e. the number of siblings.
func (nd *Node) Degree() int {
	return len(nd.siblings)
}

// Distance returns the distance to an end node. If the nodes are not in the
// same tree ErrNodesNotInSameTree will be returned.
func (nd *Node) Distance(node *Node) (int, error) {
	path, err := nd.Path(node)
	return len(path) - 1, err
}

// Get returns the node's stored data.
func (nd *Node) Get() interface{} {
	return nd.data
}

// Height returns the node's height, i.e. the longest downward path to a leaf.
func (nd *Node) Height() int {
	var height int
	l := nd.Level()

	f := func(node *Node, data interface{}) {
		if len(node.siblings) == 0 {
			if h := node.Level() - l; h > height {
				height = h
			}
		}
	}

	nd.Walk(nil, f)
	return height
}

// Index returns the index of a child in the node's list of siblings. If it
// cannot be found it returns ErrNodeNotFound.
func (nd *Node) Index(child *Node) (int, error) {
	for i, sbl := range nd.siblings {
		if sbl == child {
			return i, nil
		}
	}

	return -1, ErrNodeNotFound
}

// Level returns the node's level, i.e. the zero-based counting of edges along
// the path to the root node. It is the same as its depth.
func (nd *Node) Level() int {
	n := 0
	for p := nd.parent; p != nil; p = p.parent {
		n++
	}
	return n
}

// Link links a set of nodes. They will be inserted just before the node's
// sibling with the provide index. If the index is negative i will be set to
// zero. If the index is larger than the index of the node's last sibling, the
// nodes will be appended to the the node's siblings.
func (nd *Node) Link(index int, nodes ...*Node) error {
	descendants := make(map[*Node]dummyType)
	var found bool

	// test for duplicates in the provided nodes
	for _, n := range nodes {
		f := func(node *Node, data interface{}) {
			if !found {
				if _, found = descendants[node]; !found {
					descendants[node] = dummy
				}
			}
		}
		n.Walk(nil, f)

		if found {
			return ErrDuplicateNodeFound
		}
	}

	// tests for duplicates in the node's tree
	root := nd.Root()
	f := func(node *Node, data interface{}) {
		if !found {
			_, found = descendants[node]
		}
	}

	root.Walk(nil, f)
	if found {
		return ErrDuplicateNodeFound
	}

	for _, n := range nodes {
		n.parent = nd
	}

	nd.siblings = insertNodes(nd.siblings, nodes, index)

	return nil
}

// New returns a new node with some data stored into it.
func New(data interface{}) *Node {
	return &Node{data: data, siblings: make([]*Node, 0)}
}

// Parent returns the node's parent. When the parent doesn't exist
// ErrParentMissing will be returned.
func (nd *Node) Parent() (*Node, error) {
	parent := nd.parent
	if parent == nil {
		return nil, ErrParentMissing
	}
	return parent, nil
}

// Path returns the path to a node. The node and the end node are included into
// the result. If the nodes are not in the same tree, ErrNodesNotInSameTree will
// be returned.
func (nd *Node) Path(node *Node) ([]*Node, error) {
	up := append([]*Node{nd}, nd.Ancestors()...)
	down := append([]*Node{node}, node.Ancestors()...)
	return mergePaths(up, down)
}

// RemoveAllSiblings removes all the node's siblings. It returns a slice with
// the removed siblings. Their parents are invalidated.
func (nd *Node) RemoveAllSiblings() []*Node {
	sblngs := nd.siblings
	nd.siblings = []*Node{}

	for _, n := range sblngs {
		n.parent = nil
	}

	return sblngs
}

// RemoveSibling removes the node's sibling with the provided index. It returns
// the removed sibling. Its parent is invalidated. If there is no node with the
// given index, ErrNodeNotFound will be returned.
func (nd *Node) RemoveSibling(index int) (*Node, error) {
	l := len(nd.siblings)
	if index < 0 || index >= l {
		return nil, ErrNodeNotFound
	}

	siblings := make([]*Node, l-1)
	copy(siblings, nd.siblings[:index])
	node := nd.siblings[index]
	copy(siblings[index:], nd.siblings[index+1:])
	nd.siblings = siblings

	node.parent = nil
	return node, nil
}

// Root returns the node's root.
func (nd *Node) Root() *Node {
	node := nd
	for node.parent != nil {
		node = node.parent
	}
	return node
}

// Set stores the node's data.
func (nd *Node) Set(data interface{}) {
	nd.data = data
}

// String creates a string that displays the node's content and recursivly the
// contents of all of its descendants.
func (nd *Node) String() string {
	sb := strings.Builder{}

	fmt.Fprintf(&sb, "%v", nd.data)
	if len(nd.siblings) > 0 {
		fmt.Fprintf(&sb, "[")
		space := ""
		for _, sbl := range nd.siblings {
			fmt.Fprintf(&sb, "%s%s", space, sbl.String())
			space = " "
		}
		fmt.Fprintf(&sb, "]")
	}
	return sb.String()
}

// Walk executes f for the node and all of its descendants
func (nd *Node) Walk(data interface{}, f WalkFunc) {
	f(nd, data)
	for _, sbl := range nd.siblings {
		sbl.Walk(data, f)
	}
}
