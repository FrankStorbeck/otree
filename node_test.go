package otree

import (
	"errors"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	s0 := "s0"
	nd := New(s0)
	if got, ok := nd.Data.(string); !ok || got != s0 {
		if !ok {
			t.Errorf("New(%q) doesn't hold a strinq", s0)
		} else {
			t.Errorf("New(%q) returns data %q, should be %q", s0, got, s0)
		}
	}
}

func TestParent(t *testing.T) {
	root := New("root")
	s0 := New("s0")
	root.Link(AtStart, s0)

	tests := []struct {
		node, parent *Node
		err          error
	}{
		{s0, root, nil},
		{root, nil, ErrParentMissing},
	}
	for i, tst := range tests {
		got, err := tst.node.Parent()
		switch {
		case err != nil && tst.err == nil:
			t.Errorf("%d: Parent() returns an error %q, should be nil",
				i, err.Error())
		case err == nil && tst.err != nil:
			t.Errorf("%d: Parent() returns no error, should be %q",
				i, tst.err.Error())
		case err != nil && tst.err != nil && err != tst.err:
			t.Errorf("%d: Parent() returns error %q, should be %q",
				i, err.Error(), tst.err.Error())
		case err == nil && tst.err == nil:
			if got != tst.parent {
				t.Errorf("%d: LinkChildren() returns %q, should be %q",
					i, got.Data, tst.parent.Data)
			}
		}
	}
}

func TestString(t *testing.T) {
	root := New("root")
	s0 := New("s0")
	root.Link(AtStart, s0)

	got := root.String()
	want := "root[s0]"
	if got != want {
		t.Errorf("String() returns %q, should be %q", got, want)
	}
}

func TestLink(t *testing.T) {
	root := New("root")
	s0 := New("s0")
	parent := New("parent")
	child := New("child")
	otherParent := New("otherParent")

	tests := []struct {
		root   *Node
		where  int
		nodes  []*Node
		err    error
		want   string
		degree int
	}{
		{root, 0, []*Node{s0}, nil, "root[s0]", 1},
		{root, AtStart, []*Node{New("s1"), New("s2")}, nil, "root[s1 s2 s0]", 3},
		{root, 1, []*Node{New("s3"), New("s4")}, nil, "root[s1 s3 s4 s2 s0]", 5},
		{root, AtEnd, []*Node{New("s5")}, nil, "root[s1 s3 s4 s2 s0 s5]", 6},
		{root, AtStart, []*Node{root}, ErrDuplicateNodeFound, "", 6},
		{parent, AtEnd, []*Node{child}, nil, "parent[child]", 1},
		{parent, AtEnd, []*Node{parent}, ErrDuplicateNodeFound, "", 1},
		{parent, AtEnd, []*Node{child}, ErrDuplicateNodeFound, "", 1},
		{root, AtEnd, []*Node{parent, parent}, ErrDuplicateNodeFound, "", 6},
		{root, AtEnd, []*Node{parent}, nil, "root[s1 s3 s4 s2 s0 s5 parent[child]]", 7},
		{root, AtEnd, []*Node{child}, ErrDuplicateNodeFound, "", 7},
		{otherParent, AtEnd, []*Node{child}, nil, "otherParent[child]", 1},
		{root, AtEnd, []*Node{otherParent}, ErrDuplicateNodeFound, "", 7},
	}
	for i, tst := range tests {
		var got string
		err := tst.root.Link(tst.where, tst.nodes...)
		if err == nil {
			got = tst.root.String()
		}

		switch {
		case err != nil && tst.err == nil:
			t.Errorf("%d: Link() returns an error %q, should be nil",
				i, err.Error())
		case err == nil && tst.err != nil:
			t.Errorf("%d: Link() returns no error, should be %q",
				i, tst.err.Error())
		case err != nil && tst.err != nil && err != tst.err:
			t.Errorf("%d: Link() returns error %q, should be %q",
				i, err.Error(), tst.err.Error())
		case err == nil && tst.err == nil:
			if got != tst.want {
				t.Errorf("%d: Link() returns %q, should be %q",
					i, got, tst.want)
			}
		}
	}
}

func TestDegree(t *testing.T) {
	root := New("root")

	tests := []struct {
		root   *Node
		nodes  []*Node
		degree int
	}{
		{root, []*Node{New("s0")}, 1},
		{root, []*Node{New("s1"), New("s2")}, 3},
	}

	for i, tst := range tests {
		tst.root.Link(0, tst.nodes...)
		if d := tst.root.Degree(); d != tst.degree {
			t.Errorf("%d: Degree() returns %d, should be %d",
				i, d, tst.degree)
		}
	}
}

func TestRoot(t *testing.T) {
	root := New("root")
	s0 := New("s0")
	root.Link(AtStart, s0)
	if got := s0.Root(); got != root {
		t.Errorf("Root() returns %q, should be %q",
			got.String(), root.String())
	}
}

func TestSiblingIndex(t *testing.T) {
	root := New("root")
	sbls := []*Node{New(0), New(1), New(2), New(3), New(4)}
	root.Link(AtStart, sbls...)
	for i, sbl := range sbls {
		idx, err := root.SiblingIndex(sbl)
		if err != nil {
			t.Errorf("SiblingIndex(sbl) returns error %q, should be nil",
				err.Error())
		} else if i != idx {
			t.Errorf("SiblingIndex(sbl) returns %d, should be %d", idx, i)
		}
	}

	nd := New(5)
	_, err := root.SiblingIndex(nd)
	if err == nil {
		t.Errorf("SiblingIndex(%q) returns no error, should be %q",
			nd.String(), ErrNodeNotFound.Error())
	} else if !errors.Is(err, ErrNodeNotFound) {
		t.Errorf("SiblingIndex(nd) returns error %q, should be %q",
			err.Error(), ErrNodeNotFound.Error())
	}
}

func TestSibling(t *testing.T) {
	root := New("root")
	sbls := []*Node{New(0), New(1), New(2), New(3), New(4)}
	root.Link(AtStart, sbls...)

	for i, sbl := range root.Siblings() {
		nd, err := root.Sibling(i)
		if err != nil {
			t.Errorf("root.Sibling(%d) returns error %q, should be nil", i, err.Error())
		} else if nd != sbl {
			t.Errorf("root.Sibling(%d) returns %q, should be %q", i, nd.String(), sbl.String())
		}
	}

	_, err := root.Sibling(5)
	if err == nil {
		t.Errorf("root.Sibling(5) returns no error, should be %q", ErrNodeNotFound)
	}
}

func TestIndex(t *testing.T) {
	root := New("root")
	sbls := []*Node{New(0), New(1), New(2), New(3), New(4)}
	root.Link(AtStart, sbls...)
	for i, sbl := range sbls {
		idx, err := sbl.Index()
		if err != nil {
			t.Errorf("%q.Index() returns error %q, should be nil",
				sbl.String(), err.Error())
		} else if i != idx {
			t.Errorf("%q.Index() returns %d, should be %d",
				sbl.String(), idx, i)

		}
	}

	nd := New(5)
	_, err := nd.Index()
	if err == nil {
		t.Errorf("%q.Index() returns no error, should be %q",
			nd.String(), ErrNodeNotFound.Error())
	} else if !errors.Is(err, ErrParentMissing) {
		t.Errorf("%q.Index() returns error %q, should be %q",
			nd.String(), err.Error(), ErrParentMissing.Error())
	}
}

func TestLevel(t *testing.T) {
	root := New("root")
	child := New(0)
	grandChild := New(1)

	root.Link(0, child)
	child.Link(0, grandChild)

	if got := root.Level(); got != 0 {
		t.Errorf("root.Level() returns %d, should be 0", got)
	}
	if got := child.Level(); got != 1 {
		t.Errorf("child.Level() returns %d, should be 1", got)
	}
	if got := grandChild.Level(); got != 2 {
		t.Errorf("grandChild.Level() returns %d, should be 2", got)
	}
}

func TestAncestors(t *testing.T) {
	root := New("root")
	child := New(0)
	grandChild := New(1)

	root.Link(0, child)
	child.Link(0, grandChild)

	a := grandChild.Ancestors()
	if got := len(a); got != 2 {
		t.Errorf("len(grandChild.Ancestors()) returns %d, should be 2", got)
	}
	if a[0] != child {
		t.Errorf("grandChild.Ancestors()[0] is not child")
	}
	if a[1] != root {
		t.Errorf("grandChild.Ancestors()[1] is not root")
	}
}

func TestHeight(t *testing.T) {
	root := New("root")
	child := New(0)
	grandChild := New(1)

	root.Link(0, child)
	child.Link(0, grandChild)

	if got := root.Height(); got != 2 {
		t.Errorf("root.Height() returns %d, should be 2", got)
	}
	if got := child.Height(); got != 1 {
		t.Errorf("child.Height() returns %d, should be 1", got)
	}
	if got := grandChild.Height(); got != 0 {
		t.Errorf("grandChild.Height() returns %d, should be 0", got)
	}
}

func TestPathAndDistance(t *testing.T) {
	root := New(0)
	children := []*Node{New(10), New(11), New(12)}
	grandChildren1 := []*Node{New(20), New(21), New(22)}
	greatGrandChildren1 := []*Node{New(30), New(31), New(32)}

	root.Link(0, children...)
	children[0].Link(0, grandChildren1...)
	grandChildren1[0].Link(0, greatGrandChildren1...)

	tests := []struct {
		start, end *Node
		want       string
		err        error
		distance   int
	}{
		{grandChildren1[2], greatGrandChildren1[0], "[22 10 20 30]", nil, 3},
		{greatGrandChildren1[0], grandChildren1[2], "[30 20 10 22]", nil, 3},
		{greatGrandChildren1[0], children[0], "[30 20 10]", nil, 2},
		{greatGrandChildren1[0], children[1], "[30 20 10 0 11]", nil, 4},
		{children[0], children[0], "[10]", nil, 0},
		{children[0], New(-1), "[]", ErrNodesNotInSameTree, -1},
		{New(-1), children[0], "[]", ErrNodesNotInSameTree, -1},
	}

	for i, tst := range tests {
		path, err := tst.start.Path(tst.end)
		switch {
		case err != nil && tst.err == nil:
			t.Errorf("%d: Path() returns an error %q, should be nil",
				i, err.Error())
		case err == nil && tst.err != nil:
			t.Errorf("%d: Path() returns no error, should be %q",
				i, tst.err.Error())
		case err != nil && tst.err != nil && err != tst.err:
			t.Errorf("%d: Path() returns error %q, should be %q",
				i, err.Error(), tst.err.Error())
		case err == nil && tst.err == nil:
			got := "["
			space := ""
			for _, nd := range path {
				got += fmt.Sprintf("%s%d", space, nd.Data)
				space = " "
			}
			got += "]"

			if got != tst.want {
				t.Errorf("Path() results in %q, should be %q", got, tst.want)
			}
			if d, _ := tst.start.Distance(tst.end); d != tst.distance {
				t.Errorf("Distance() is %d, should be %d", d, tst.distance)
			}
		}
	}
}

func TestRemoveAllSiblings(t *testing.T) {
	root := New("root")
	children := []*Node{New("s0"), New("s1")}
	root.Link(0, children...)

	got := root.RemoveAllSiblings()
	if root.Degree() != 0 {
		t.Errorf("RemoveAllSiblings() failed")
	} else {
		different := false
		for i, n := range children {
			if n != got[i] {
				different = true
			}
		}
		if different {
			t.Errorf("RemoveAllSiblings() returns %v, should be %v", got, children)
		}
	}
}

func TestRemoveSibling(t *testing.T) {
	root := New("root")
	children := []*Node{New("s0"), New("s1"), New("s2")}
	root.Link(0, children...)

	tests := []struct {
		i    int
		want *Node
		err  error
	}{
		{0, children[0], nil},
		{1, children[2], nil},
		{-1, children[2], ErrNodeNotFound},
		{1, children[2], ErrNodeNotFound},
		{0, children[1], nil},
	}

	for _, tst := range tests {
		got, err := root.RemoveSibling(tst.i)
		switch {
		case err != nil && tst.err == nil:
			t.Errorf("RemoveSibling(%d) returns an error %q, should be nil",
				tst.i, err.Error())
		case err == nil && tst.err != nil:
			t.Errorf("RemoveSibling(%d) returns no error, should be %q",
				tst.i, tst.err.Error())
		case err != nil && tst.err != nil && err != tst.err:
			t.Errorf("RemoveSibling(%d) returns error %q, should be %q",
				tst.i, err.Error(), tst.err.Error())
		case err == nil && tst.err == nil:
			if got != tst.want {
				t.Errorf("RemoveSibling(%d) returns error %q, should be %q",
					tst.i, got.String(), tst.want.String())
			}
		}
	}
}

func TestReplaceSibling(t *testing.T) {
	root := New("root")
	root.Link(0, New("s0"), New("s1"), New("s2"))

	tests := []struct {
		i     int
		nodes []*Node
		err   error
		want  string
	}{
		{1, []*Node{New("sa")}, nil, "root[s0 sa s2]"},
		{1, []*Node{New("sb"), New("sc")}, nil, "root[s0 sb sc s2]"},
		{0, []*Node{New("sd"), New("se")}, nil, "root[sd se sb sc s2]"},
		{4, []*Node{New("sf")}, nil, "root[sd se sb sc sf]"},
		{AtStart, []*Node{New("sg")}, ErrNodeNotFound, ""},
		{5, []*Node{New("sh")}, ErrNodeNotFound, ""},
	}

	for _, tst := range tests {
		_, err := root.ReplaceSibling(tst.i, tst.nodes...)
		switch {
		case err != nil && tst.err == nil:
			t.Errorf("RemoveSibling(%d) returns an error %q, should be nil",
				tst.i, err.Error())
		case err == nil && tst.err != nil:
			t.Errorf("RemoveSibling(%d) returns no error, should be %q",
				tst.i, tst.err.Error())
		case err != nil && tst.err != nil && err != tst.err:
			t.Errorf("RemoveSibling(%d) returns error %q, should be %q",
				tst.i, err.Error(), tst.err.Error())
		case err == nil && tst.err == nil:
			if got := root.String(); got != tst.want {
				t.Errorf("ReplaceSibling(%d, %v) returns error %q, should be %q",
					tst.i, tst.nodes, got, tst.want)
			}
		}
	}
}

func TestIsLeaf(t *testing.T) {
	root := New("root")
	leaf := New("s0")
	root.Link(0, leaf)

	if root.IsLeaf() {
		t.Errorf("root.IsLeaf() returns true, should be false.")
	}
	if !leaf.IsLeaf() {
		t.Errorf("leaf.IsLeaf() returns false, should be true.")
	}
}

func TestReplace(t *testing.T) {
	root := New("root")
	s0 := New("s0")
	s1 := New("s1")
	s2 := New("s2")
	root.Link(0, s0, s1, s2)

	tests := []struct {
		node  *Node
		nodes []*Node
		err   error
		want  string
	}{
		{s1, []*Node{New("sa")}, nil, "root[s0 sa s2]"},
		{s0, []*Node{New("sb"), New("sc")}, nil, "root[sb sc sa s2]"},
		{root, []*Node{New("sd"), New("se")}, ErrCannotReplaceRootNode, ""},
	}

	for _, tst := range tests {
		_, err := tst.node.Replace(tst.nodes...)
		if err != nil {
			println(err.Error())
		}
		switch {
		case err != nil && tst.err == nil:
			t.Errorf("Replace() returns an error %q, should be nil",
				err.Error())
		case err == nil && tst.err != nil:
			t.Errorf("Replace() returns no error, should be %q",
				tst.err.Error())
		case err != nil && tst.err != nil:
			if err != tst.err {
				t.Errorf("Replace() returns error %q, should be %q",
					err.Error(), tst.err.Error())
			}
		default:
			if s := root.String(); s != tst.want {
				t.Errorf("Replace() results in %q, should be %q", s, tst.want)
			}
		}

	}
}
