// Package otree implements an ordered tree  structure.
// See: https://en.wikipedia.org/wiki/Tree_(data_structure)
package otree

// Breadth returns the breadth (i.e. the number of leaves) of the tree starting
// at the root of node. If sub is true it returns the size of the subtree for
// which node is the root.
func Breadth(node *Node, sub bool) int {
	breadth := 0

	f := func(node *Node, data interface{}) {
		if node.IsLeaf() {
			breadth++
		}
	}

	selectRoot(node, sub).Walk(f, nil)
	return breadth
}

// Degree returns the degree (i.e. the maximum degree of all nodes) of the tree
// starting at the root of node. If sub is true it returns the size of the
// subtree starting at node.
func Degree(node *Node, sub bool) int {
	degree := 0

	f := func(node *Node, data interface{}) {
		if d := node.Degree(); d > degree {
			degree = d
		}
	}

	selectRoot(node, sub).Walk(f, nil)
	return degree
}

// Height returns the height of the tree starting at the root of node. If sub
// is true it returns the height of the subtree starting at node.
func Height(node *Node, sub bool) int {
	return selectRoot(node, sub).Height()
}

// Size returns the size (i.e the number of nodes) of the tree starting at the
// root of node. If sub is true it returns the size of the subtree starting at
// node.
func Size(node *Node, sub bool) int {
	n := 0

	f := func(nd *Node, data interface{}) {
		n++
	}
	selectRoot(node, sub).Walk(f, nil)

	return n
}

// Width returns the width (i.e. the number of nodes) for a given level of a
// tree starting at the root of node. If sub is true it returns the width of the
// subtree starting at node.
func Width(node *Node, level int, sub bool) int {
	width := 0

	f := func(node *Node, data interface{}) {
		if node.Level() == level {
			width++
		}
	}

	selectRoot(node, sub).Walk(f, nil)
	return width
}
