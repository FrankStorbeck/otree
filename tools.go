package otree

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
