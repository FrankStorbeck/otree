package otree

// Node represents internal and external nodes. Internal nodes have children,
// external nodes store data.
type Node struct {
	data     interface{} // data
	level    int         // the level of the node
	parent   *Node       // parent node
	siblings []*Node     // sibling nodes
}

// WalkFunc is a function type that can be performed on nodes.
type WalkFunc func(node *Node, data interface{})

// Ancestors returns the node's ancestors.The first one is it parent, the next
// one its grandparent and so on.
func (nd *Node) Ancestors() []*Node {
	ancstrs := []*Node{}

	for p := nd.parent; p != nil; p = p.parent {
		ancstrs = append(ancstrs, p)
	}
	return ancstrs
}

// Distance returns the distance to an end node. The distance to itself is zero,
// The distance to non reachable node is -1.
func (nd *Node) Distance(endNode *Node) int {
	return len(nd.Path(endNode)) - 1
}

// NewNode returns a new node
func NewNode(data interface{}) *Node {
	return &Node{data: data, siblings: make([]*Node, 0)}
}

// Get returns the stored data.
func (nd *Node) Get() interface{} {
	return nd.data
}

// Height returns the height
func (nd *Node) Height() int {
	var height int
	f := func(node *Node, data interface{}) {
		if len(node.siblings) == 0 {
			if h := node.level - nd.level; h > height {
				height = h
			}
		}
	}
	nd.Walk(nil, f)
	return height
}

// Index returns the index of a child in the list of siblings. If it cannot be
// found it returns ErrNoNodeFound.
func (nd *Node) Index(child *Node) (int, error) {
	for i, sbl := range nd.siblings {
		if sbl == child {
			return i, nil
		}
	}

	return -1, ErrNoNodeFound
}

// Level returns the level of the node. It is the same as its depth.
func (nd *Node) Level() int {
	return nd.level
}

// Parent returns the parent node. When the parent doesn't exist
// ErrParentMissing will be returned.
func (nd *Node) Parent() (parent *Node, err error) {
	parent = nd.parent
	if parent == nil {
		err = ErrParentMissing
	}
	return
}

// Path returns the path to a node. The node and the end node are included into
// the result.
func (nd *Node) Path(endNode *Node) []*Node {
	up := append([]*Node{nd}, nd.Ancestors()...)
	down := append([]*Node{endNode}, endNode.Ancestors()...)
	return mergePaths(up, down)
}

// Set stores the data
func (nd *Node) Set(data interface{}) {
	nd.data = data
}

// Walk executes f for nd and all its descendants
func (nd *Node) Walk(data interface{}, f WalkFunc) {
	f(nd, data)
	for _, sbl := range nd.siblings {
		sbl.Walk(data, f)
	}
}
