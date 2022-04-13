// Package otree implements an ordered tree  structure.
// See: https://en.wikipedia.org/wiki/Tree_(data_structure)
package otree

type noneT struct{}

var none = noneT{}

// Tree represents an ordered tree.
type Tree struct {
	root    *Node           // root
	present map[*Node]noneT // map with all nodes present in the tree
}

// Breadth returns the breadth
func (tr *Tree) Breadth() int {
	breadth := 0

	f := func(node *Node, data interface{}) {
		if node.Degree() == 0 {
			breadth++
		}
	}

	tr.root.Walk(nil, f)
	return breadth
}

// Degree returns the degree of the tree
func (tr *Tree) Degree() int {
	degree := 0

	f := func(node *Node, data interface{}) {
		if d := node.Degree(); d > degree {
			degree = d
		}
	}

	tr.root.Walk(nil, f)
	return degree
}

// Height returns the height of the tree
func (tr *Tree) Height() int {
	return tr.root.Height()
}

// New returns a new initialised Tree with data stored in the root node
func New(data interface{}) Tree {
	tr := Tree{
		root:    NewNode(data),
		present: make(map[*Node]noneT, 0),
	}
	tr.present[tr.root] = none
	return tr
}

// LinkChildren links a set of children to the siblings of an existing parent
// node in the tree. If the parent is not a node in the tree ErrNodeNotFound wil
// be returned. The children will be inserted just before the sibling with index
// i. However if i is negative they wil be inserted before the first sibling. If
// i is larger than the index of the last sibling they will be appended to the
// end of the siblings. The children must not have children them selves.
// Otherwise ErrNodeMustNotHaveSiblings will be returned.
func (tr *Tree) LinkChildren(parent *Node, i int, children ...*Node) error {
	if _, found := tr.present[parent]; !found {
		return ErrNodeNotFound
	}

	for _, nd := range children {
		if _, found := tr.present[nd]; found {
			return ErrDuplicateNodeFound
		}
		if len(nd.siblings) > 0 {
			return ErrNodeMustNotHaveSiblings
		}
	}

	level := parent.level + 1
	for _, nd := range children {
		nd.parent = parent
		nd.level = level
		tr.present[nd] = none // mark nodes as present in the tree
	}

	parent.siblings = insertNodes(parent.siblings, children, i)
	return nil
}

// RemoveNode removes node and all its descendants from the tree.
func (tr Tree) RemoveNode(node *Node) error {
	if node == nil {
		return ErrNoNodeFound
	}
	if node == tr.root {
		return ErrCannotRemoveRootNode
	}
	p, err := node.Parent()
	if err != nil {
		return err
	}

	i, err := p.Index(node)
	if err != nil {
		return err
	}
	_, err = tr.RemoveSibling(p, i)

	return err
}

// RemoveSibling removes the sibling with index i from a parent. It returns a
// pointer to the removed sibling.
func (tr Tree) RemoveSibling(parent *Node, i int) (*Node, error) {
	l := len(parent.siblings)

	siblings := make([]*Node, l-1)
	copy(siblings, parent.siblings[:i])
	deletedNode := parent.siblings[i]
	copy(siblings[i:], parent.siblings[i+1:])
	parent.siblings = siblings

	deletedNode.Walk(nil, func(nd *Node, data interface{}) {
		delete(tr.present, nd)
	})

	return deletedNode, nil
}

// RemoveSiblings removes all siblings from a parent.
func (tr Tree) RemoveSiblings(parent *Node) {
	for _, sblng := range parent.siblings {
		sblng.Walk(nil, func(nd *Node, data interface{}) {
			delete(tr.present, nd)
		})
	}
	parent.siblings = []*Node{}
}

// Root returns the root node
func (tr *Tree) Root() *Node {
	return tr.root
}

// Size returns the size of the tree
func (tr *Tree) Size() int {
	return len(tr.present)
}

// String creates a string that displays the content of a tree
func (tr *Tree) String() string {
	return tr.root.String()
}

// Width returns the width for a given level
func (tr *Tree) Width(level int) int {
	width := 0

	f := func(node *Node, data interface{}) {
		if node.Level() == level {
			width++
		}
	}

	tr.root.Walk(nil, f)
	return width
}
