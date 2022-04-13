package otree

// markAndAssignLevel marks the node and its descendants as present in the tree
// and reassigns the levels for them.
func (tr *Tree) mark(nd *Node) {
	tr.present[nd] = none // mark nodes as present
	for _, sblng := range nd.siblings {
		tr.mark(sblng)
	}
}

// insertNodes inserts nodes2 into nodes1 before index i
func insertNodes(nodes1, nodes2 []*Node, i int) []*Node {
	l1 := len(nodes1)
	l2 := len(nodes2)

	if i > l1 {
		i = l1
	} else if i < 0 {
		i = 0
	}

	r := make([]*Node, l1+l2)
	copy(r, nodes1[:i])
	copy(r[i:], nodes2)
	copy(r[i+l2:], nodes1[i:])
	return r
}

// invertSlice inverts the sequence of the nodes
func invertSlice(nodes []*Node) []*Node {
	for i, j := 0, len(nodes)-1; i < j; i, j = i+1, j-1 {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	}
	return nodes
}

// mergePaths merges the up and down paths via the lowest shared node
func mergePaths(up, down []*Node) []*Node {
	up, down = invertSlice(up), invertSlice(down)

	var l int
	if l = len(up); l > len(down) {
		l = len(down)
	}

	if l == 0 || up[0] != down[0] {
		return []*Node{}
	}

	i := 1
	common := up[0]
	for i < l {
		if up[i] != down[i] {
			break
		}
		common = up[i]
		i++
	}

	r := make([]*Node, (len(up)-i)+(len(down)-i)+1)
	copy(r, invertSlice(up[i:]))
	l = len(up[i:])
	r[l] = common
	copy(r[l+1:], down[i:])
	return r
}
