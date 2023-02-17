package otree

import (
	"fmt"
	"strings"
)

// Node represents internal and external/leaf nodes. It can hold links to other
// nodes which are called its children. These children are named the node's
// siblings. All nodes other than the root node hold one parent node.
// Internal nodes have children, external nodes don't.
type Node struct {
	data     interface{} // data
	parent   *Node       // parent node
	siblings []*Node     // sibling nodes
}

// WalkFunc is a function type that can be performed on all nodes in a
// (sub)tree. Walk() sets nd to the node for wich the function must be called.
type WalkFunc func(nd *Node, data interface{})

type dummyType struct{}

var dummy = dummyType{}

// Ancestors returns the node's ancestors. The first one is it parent, the next
// one its grandparent and so on until the root node is found.
func (nd *Node) Ancestors() []*Node {
	ancestors := []*Node{}

	for p := nd.parent; p != nil; p = p.parent {
		ancestors = append(ancestors, p)
	}
	return ancestors
}

// Degree returns nd's degree, i.e. the number of siblings.
func (nd *Node) Degree() int {
	return len(nd.siblings)
}

// Distance returns nd's distance (the number of edges) to node. If nd
// and node are not in the same tree ErrNodesNotInSameTree will be
// returned.
func (nd *Node) Distance(node *Node) (int, error) {
	path, err := nd.Path(node)
	return len(path) - 1, err
}

// Height returns nd's height, i.e. the longest downward path to a leaf.
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

	nd.Walk(f, nil)
	return height
}

// SiblingIndex returns the index of child in nd's list of siblings. If it
// cannot be found it returns ErrNodeNotFound.
func (nd *Node) SiblingIndex(child *Node) (int, error) {
	for i, sbl := range nd.siblings {
		if sbl == child {
			return i, nil
		}
	}

	return -1, ErrNodeNotFound
}

// Index returns the index in the list of siblings to which nd belongs.
// Finding the index of a root node results in ErrParentMissing to be
// returned.
func (nd *Node) Index() (int, error) {
	p, err := nd.Parent()
	if err != nil {
		return -1, err
	}
	return p.SiblingIndex(nd)
}

// IsLeaf tells if nd is an external/leaf node.
func (nd *Node) IsLeaf() bool {
	return len(nd.siblings) == 0
}

// Level returns nd's level, i.e. the zero-based counting of edges along
// the path to the root node. It is the same as its depth.
func (nd *Node) Level() int {
	n := 0
	for p := nd.parent; p != nil; p = p.parent {
		n++
	}
	return n
}

// Link links nodes to nd. They will be inserted just before the
// child in the list of siblings with index. If index is negative it will be
// set to zero. If the index is larger than the index of the nd's last sibling
// nodes will be appended to the the node's siblings.
// The constants AtStart and AtEnd can be used to link nodes at the start or
// end of the list of siblings.
func (nd *Node) Link(index int, nodes ...*Node) error {
	newNodes := make(map[*Node]dummyType)
	found, collectNodes := false, true

	// f is a WalkFunc to test if there are any duplicate nodes in a tree
	f := func(node *Node, data interface{}) {
		if !found {
			_, found = newNodes[node]
			if !found && collectNodes {
				newNodes[node] = dummy
			}
		}
	}

	// test for duplicates in the provided nodes
	for _, n := range nodes {
		n.Walk(f, nil)
		if found {
			return ErrDuplicateNodeFound
		}
	}

	// tests for duplicates in the node's tree
	collectNodes = false
	nd.Root().Walk(f, nil)
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

// Parent returns nd's parent. When the parent doesn't exist ErrParentMissing
// will be returned.
func (nd *Node) Parent() (*Node, error) {
	parent := nd.parent
	if parent == nil {
		return nil, ErrParentMissing
	}
	return parent, nil
}

// Path returns nd's path to node. nd and node are included into the result.
// If the nodes are not in the same tree, ErrNodesNotInSameTree will be returned.
func (nd *Node) Path(node *Node) ([]*Node, error) {
	up := append([]*Node{nd}, nd.Ancestors()...)
	down := append([]*Node{node}, node.Ancestors()...)
	return mergePaths(up, down)
}

// RemoveAllSiblings removes all nd's siblings. It returns a slice with
// the removed siblings. Their parents are invalidated.
func (nd *Node) RemoveAllSiblings() []*Node {
	sblngs := nd.siblings
	nd.siblings = []*Node{}

	for _, n := range sblngs {
		n.parent = nil
	}

	return sblngs
}

// RemoveSibling removes the nd's child with the provided index in the list
// of siblings. It returns the removed sibling. Its parents is invalidated.
// If there is no node with the given index, ErrNodeNotFound will be returned.
func (nd *Node) RemoveSibling(index int) (*Node, error) {
	l := nd.Degree()
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

// Replace replaces nd by nodes. It returns itself with an invalidated
// parent. if nd is the root node ErrCannotReplaceRootNode will be returned.
func (nd *Node) Replace(nodes ...*Node) (*Node, error) {
	p, err := nd.Parent()
	if err != nil {
		return nd, ErrCannotReplaceRootNode
	}

	i, err := nd.Index()
	if err != nil {
		return nd, err
	}

	return p.ReplaceSibling(i, nodes...)
}

// ReplaceSibling replaces the child in the list of sibling with the provided
// index by nodes. It returns the replaced sibling with an invalidated parent.
func (nd *Node) ReplaceSibling(index int, nodes ...*Node) (*Node, error) {
	if index < 0 || index >= nd.Degree() {
		return nd, ErrNodeNotFound
	}
	node, err := nd.RemoveSibling(index)
	if err != nil {
		return node, err
	}
	err = nd.Link(index, nodes...)
	return node, err
}

// Root returns the root of the tree to which nd belongs.
func (nd *Node) Root() *Node {
	node := nd
	for node.parent != nil {
		node = node.parent
	}
	return node
}

// SetData stores the node's data.
func (nd *Node) SetData(data interface{}) {
	nd.data = data
}

// Sibling returns nd's child in the list of siblings with the provided
// index.
func (nd *Node) Sibling(index int) (*Node, error) {
	if index < 0 || index >= len(nd.siblings) {
		return nil, ErrNodeNotFound
	}
	return nd.siblings[index], nil
}

// Siblings returns all the siblings
func (nd *Node) Siblings() []*Node {
	return nd.siblings
}

// String creates a string that displays nd's content and recursivly the
// contents of all of its newNodes.
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

// Walk executes f for nd and all of its newNodes.
func (nd *Node) Walk(f WalkFunc, data interface{}) {
	f(nd, data)
	for _, sbl := range nd.siblings {
		sbl.Walk(f, data)
	}
}
