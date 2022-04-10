package otree

// Node represents internal and external nodes. Internal nodes have children,
// external nodes store data.
type Node struct {
	data     interface{} // data
	parent   *Node       // parent node
	siblings []*Node     // sibling nodes
}

// NewNode returns a new node
func NewNode(data interface{}) *Node {
	return &Node{data: data, siblings: make([]*Node, 0)}
}

// Get returns the stored data.
func (nd *Node) Get() interface{} {
	return nd.data
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

// Parent returns the parent node. When the parent doesn't exist
// ErrParentMissing will be returned.
func (nd *Node) Parent() (parent *Node, err error) {
	parent = nd.parent
	if parent == nil {
		err = ErrParentMissing
	}
	return
}

// Set stores the data
func (nd *Node) Set(data interface{}) {
	nd.data = data
}
