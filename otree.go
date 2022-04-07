// Package tree implements an ordered tree  structure.
// See: https://en.wikipedia.org/wiki/Tree_(data_structure)
package otree

type noneT struct{}

var none = noneT{}

// Tree represents an ordered tree.
type Tree struct {
	root       *Node           // root node
	validNodes map[*Node]noneT // map with all nodes present
}

// New returns a new initialised Tree
func New() Tree {
	return Tree{
		root:       NewNode(none),
		validNodes: make(map[*Node]noneT, 0),
	}
}
