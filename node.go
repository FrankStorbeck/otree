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

// Set stores the data
func (nd *Node) Set(data interface{}) {
	nd.data = data
}
