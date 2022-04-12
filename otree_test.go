package otree

import (
	"fmt"
	"strings"
	"testing"
)

func TestLinkNodesAndDegree(t *testing.T) {
	tr := New()
	s0 := NewNode("s0")
	nWC := NewNode("node with child")
	nWC.siblings = []*Node{NewNode("child")}

	tests := []struct {
		node   *Node
		where  int
		nodes  []*Node
		err    error
		want   string
		degree int
	}{
		{tr.root, 0, []*Node{s0}, nil, "root[s0]", 1},
		{tr.root, -1, []*Node{NewNode("s1"), NewNode("s2")}, nil, "root[s1 s2 s0]", 3},
		{tr.root, 1, []*Node{NewNode("s3"), NewNode("s4")}, nil, "root[s1 s3 s4 s2 s0]", 5},
		{tr.root, 500, []*Node{NewNode("s5")}, nil, "root[s1 s3 s4 s2 s0 s5]", 6},
		{tr.root, -1, []*Node{tr.root}, ErrDuplicateNodeFound, "", 6},
		{NewNode(""), 1, []*Node{tr.root}, ErrNodeNotFound, "", 6},
		{tr.root, 1, []*Node{nWC}, ErrNodeMustNotHaveSiblings, "", 6},
	}

	for i, tst := range tests {
		var got string
		err := tr.LinkChildren(tst.node, tst.where, tst.nodes...)
		if err == nil {
			got = print(tr.root)
		}

		switch {
		case err != nil && tst.err == nil:
			t.Errorf("%d: LinkChildren() returns an error %q, should be nil",
				i, err.Error())
		case err == nil && tst.err != nil:
			t.Errorf("%d: LinkChildren() returns no error, should be %q",
				i, tst.err.Error())
		case err != nil && tst.err != nil && err != tst.err:
			t.Errorf("%d: LinkChildren() returns error %q, should be %q",
				i, err.Error(), tst.err.Error())
		case err == nil && tst.err == nil:
			if got != tst.want {
				t.Errorf("%d: LinkChildren() returns %q, should be %q",
					i, got, tst.want)
			}
		}
		if d := tr.Degree(); d != tst.degree {
			t.Errorf("%d: Degree() returns %d, should be %d",
				i, d, tst.degree)
		}
	}
}

// print creates a string that displays the content of a node
func print(node *Node) string {
	sb := strings.Builder{}

	switch k := node.data.(type) {
	case int:
		fmt.Fprintf(&sb, "%d", k)
	case string:
		fmt.Fprintf(&sb, "%s", k)
	default:
		fmt.Fprintf(&sb, "?")
	}
	if len(node.siblings) > 0 {
		fmt.Fprintf(&sb, "[")
		space := ""
		for _, sbl := range node.siblings {
			fmt.Fprintf(&sb, "%s%s", space, print(sbl))
			space = " "
		}
		fmt.Fprintf(&sb, "]")
	}
	return sb.String()
}

func TestTreeHeight(t *testing.T) {
	tr := New()
	child := NewNode(0)
	grandChild := NewNode(1)

	tr.LinkChildren(tr.root, 0, child)
	tr.LinkChildren(child, 0, grandChild)

	if got := tr.Height(); got != 2 {
		t.Errorf("tr.Height() returns %d, should be 2", got)
	}
}
