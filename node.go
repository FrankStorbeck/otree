// Package otree implements functions for manipulating an ordered tree
// structure.
package otree

import (
	"fmt"
	"strings"
)

// Node` is a structure which may contain data and links to other nodes.
// It has zero or more `child` nodes and exactly one link to a parent node,
// exept the root node. The child nodes of a parent node are the siblings.
// These siblings have an order.
// An internal node is any node of a tree that has one or more child nodes. An
// external node, or leaf node, is any node that does not have child nodes.
type Node struct {
	Data     interface{} // stored data
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
func (nd *Node) Ancestors() (ancestors []*Node) {
	for p := nd.parent; p != nil; p = p.parent {
		ancestors = append(ancestors, p)
	}
	return
}

// Degree returns nd's degree, i.e. the number of siblings.
func (nd *Node) Degree() (d int) {
	if nd.siblings != nil {
		d = len(nd.siblings)
	}
	return
}

// Distance returns nd's distance (the number of edges) to node. If nd
// and node are not in the same tree ErrNodesNotInSameTree will be
// returned.
func (nd *Node) Distance(node *Node) (int, error) {
	path, err := nd.Path(node)
	return len(path) - 1, err
}

// Height returns nd's height, i.e. the longest downward path to a leaf.
func (nd *Node) Height() (height int) {
	l := nd.Level()

	f := func(node *Node, data interface{}) {
		if node.IsLeaf() {
			if h := node.Level() - l; h > height {
				height = h
			}
		}
	}

	nd.Walk(f, nil)
	return
}

// Index returns the index in the list of siblings to which nd belongs.
// Finding the index of a root node results in returning ErrParentMissing.
func (nd *Node) Index() (int, error) {
	p, err := nd.Parent()
	if err != nil {
		return -1, err
	}
	return p.SiblingIndex(nd)
}

// IsLeaf tells if nd is an external/leaf node.
func (nd *Node) IsLeaf() bool {
	return nd.siblings == nil || len(nd.siblings) == 0
}

// Level returns nd's level, i.e. the zero-based counting of edges along
// the path to the root node. It is the same as its depth.
func (nd *Node) Level() (n int) {
	for p := nd.parent; p != nil; p = p.parent {
		n++
	}
	return
}

// Link links nodes to nd. They will be inserted just before the
// child in the list of siblings with index. If index is negative it will be
// set to zero. If the index is larger than the index of the nd's last sibling
// nodes will be appended to the the node's siblings.
// The constants AtStart and AtEnd can be used to link nodes at the start or
// end of the list of siblings.
func (nd *Node) Link(index int, nodes ...*Node) error {
	newNodes := make(map[*Node]dummyType)
	found := false

	// f is a WalkFunc to test if there are any duplicate nodes in a tree
	f := func(node *Node, collectNodes interface{}) {
		if !found {
			_, found = newNodes[node]
			if !found && collectNodes.(bool) {
				newNodes[node] = dummy
			}
		}
	}

	// test for duplicates in the provided nodes
	for _, n := range nodes {
		if n.Walk(f, true); found {
			return ErrDuplicateNodeFound
		}
	}

	// tests for duplicates in the node's tree
	if nd.Root().Walk(f, false); found {
		return ErrDuplicateNodeFound
	}

	for _, n := range nodes {
		n.parent = nd
	}

	if nd.siblings == nil {
		nd.siblings = make([]*Node, 0)
	}
	nd.siblings = insertNodes(nd.siblings, nodes, index)

	return nil
}

// New returns a new node with some data stored into it.
func New(data interface{}) *Node {
	return &Node{Data: data}
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

// Remove removes nd from the tree. On return nd has an invalidated
// parent. The root node cannot be removed.
func (nd *Node) Remove() error {
	p := nd.parent
	if p == nil {
		return ErrCannotRemoveRootNode
	}
	i, err := nd.Index()
	if err != nil {
		return err
	}

	_, err = p.RemoveSibling(i)
	return err
}

// RemoveAllSiblings removes all nd's siblings from the tree. It returns a slice with
// the removed siblings. Their parents are invalidated.
func (nd *Node) RemoveAllSiblings() []*Node {
	if nd.IsLeaf() {
		return []*Node{}
	}
	sblngs := nd.siblings
	nd.siblings = nil

	for _, n := range sblngs {
		n.parent = nil
	}

	return sblngs
}

// RemoveSibling removes the nd's child with the provided index in the list
// of siblings. It returns the removed sibling. Its parent is invalidated.
// If there is no node with the given index, ErrNodeNotFound will be returned.
func (nd *Node) RemoveSibling(index int) (*Node, error) {
	l := nd.Degree()
	if nd.IsLeaf() || index < 0 || index >= l {
		return nil, ErrNodeNotFound
	}

	siblings := make([]*Node, l-1)
	copy(siblings, nd.siblings[:index])
	node := nd.siblings[index]
	copy(siblings[index:], nd.siblings[index+1:])

	if len(siblings) == 0 {
		nd.siblings = nil
	} else {
		nd.siblings = siblings
	}

	node.parent = nil
	return node, nil
}

// Replace replaces nd by nodes. nd's parent wil be invalidated. If nd is the
// root node ErrCannotReplaceRootNode will be returned.
func (nd *Node) Replace(nodes ...*Node) error {
	p, err := nd.Parent()
	if err != nil {
		return ErrCannotReplaceRootNode
	}

	i, err := nd.Index()
	if err != nil {
		return err
	}
	_, err = p.ReplaceSibling(i, nodes...)
	return err
}

// ReplaceSibling replaces the child in the list of siblings with the provided
// index by nodes. It returns the replaced sibling with an invalidated parent.
func (nd *Node) ReplaceSibling(index int, nodes ...*Node) (*Node, error) {
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

// Sibling returns nd's child in the list of siblings with the provided
// index.
func (nd *Node) Sibling(index int) (*Node, error) {
	if nd.IsLeaf() || index < 0 || index >= len(nd.siblings) {
		return nil, ErrNodeNotFound
	}
	return nd.siblings[index], nil
}

// Siblings returns all the siblings
func (nd *Node) Siblings() []*Node {
	if nd.IsLeaf() {
		return []*Node{}
	}
	return nd.siblings
}

// SiblingIndex returns the index of child in nd's list of siblings. If it
// cannot be found it returns ErrNodeNotFound.
func (nd *Node) SiblingIndex(child *Node) (int, error) {
	if nd.siblings != nil {
		for i, sbl := range nd.siblings {
			if sbl == child {
				return i, nil
			}
		}
	}

	return -1, ErrNodeNotFound
}

// String creates a string that displays nd's content and recursivly the
// contents of all of its newNodes.
func (nd *Node) String() string {
	sb := strings.Builder{}

	fmt.Fprintf(&sb, "%v", nd.Data)
	if nd.siblings != nil && len(nd.siblings) > 0 {
		fmt.Fprintf(&sb, "[")
		sep := ""
		for _, sbl := range nd.siblings {
			fmt.Fprintf(&sb, "%s%s", sep, sbl.String())
			sep = " "
		}
		fmt.Fprintf(&sb, "]")

	}
	return sb.String()
}

// Walk executes f for nd and all of its descendants. data will be used
// as the second argument for f.
func (nd *Node) Walk(f WalkFunc, data interface{}) {
	f(nd, data)
	if nd.siblings != nil {
		for _, sbl := range nd.siblings {
			sbl.Walk(f, data)
		}
	}
}

// WalkUp executes f for nd and all of its descendants. data will be used
// as the second argument for f. In contrast to Walk, f will be performed to
// nd's siblings before executing it on nd itself.
func (nd *Node) WalkUp(f WalkFunc, data interface{}) {
	if nd.siblings != nil {
		for _, sbl := range nd.siblings {
			sbl.Walk(f, data)
		}
	}
	f(nd, data)
}
