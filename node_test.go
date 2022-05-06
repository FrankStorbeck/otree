package otree

import (
	"errors"
	"fmt"
	"math"
	"testing"
)

func TestNew(t *testing.T) {
	s0 := "s0"
	nd := New(s0)
	if got := nd.data.(string); got != s0 {
		t.Errorf("New(%q) returns data %q, should be %q", s0, got, s0)
	}
}

func TestSetAndGet(t *testing.T) {
	tests := []struct {
		data interface{}
	}{
		{"s"},
		{8},
		{struct {
			i int
			f float32
		}{2, 3.4}},
	}
	nd := New("")
	for _, tst := range tests {
		nd.Set(tst.data)
		got := nd.Get()
		if got != tst.data {
			t.Errorf("Get() after Set(%v) returns %v, should be %v", tst.data, got, tst.data)
		}
	}
}

func TestParent(t *testing.T) {
	root := New("root")
	s0 := New("s0")
	root.Link(-1, s0)

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
					i, got.data, tst.parent.data)
			}
		}
	}
}

func TestString(t *testing.T) {
	root := New("root")
	s0 := New("s0")
	root.Link(-1, s0)

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
		{root, -1, []*Node{New("s1"), New("s2")}, nil, "root[s1 s2 s0]", 3},
		{root, 1, []*Node{New("s3"), New("s4")}, nil, "root[s1 s3 s4 s2 s0]", 5},
		{root, math.MaxInt, []*Node{New("s5")}, nil, "root[s1 s3 s4 s2 s0 s5]", 6},
		{root, -1, []*Node{root}, ErrDuplicateNodeFound, "", 6},
		{parent, math.MaxInt, []*Node{child}, nil, "parent[child]", 1},
		{parent, math.MaxInt, []*Node{parent}, ErrDuplicateNodeFound, "", 1},
		{parent, math.MaxInt, []*Node{child}, ErrDuplicateNodeFound, "", 1},
		{root, math.MaxInt, []*Node{parent, parent}, ErrDuplicateNodeFound, "", 6},
		{root, math.MaxInt, []*Node{parent}, nil, "root[s1 s3 s4 s2 s0 s5 parent[child]]", 7},
		{root, math.MaxInt, []*Node{child}, ErrDuplicateNodeFound, "", 7},
		{otherParent, math.MaxInt, []*Node{child}, nil, "otherParent[child]", 1},
		{root, math.MaxInt, []*Node{otherParent}, ErrDuplicateNodeFound, "", 7},
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
	root.Link(-1, s0)
	if got := s0.Root(); got != root {
		t.Errorf("Root() returns %q, should be %q",
			got.String(), root.String())
	}
}

func TestIndex(t *testing.T) {
	root := New("root")
	sbls := []*Node{New(0), New(1), New(2), New(3), New(4)}
	root.Link(-1, sbls...)
	for i, sbl := range sbls {
		idx, err := root.Index(sbl)
		if err != nil {
			t.Errorf("Index(sbl) returns error %q, should be nil",
				err.Error())
		} else if i != idx {
			t.Errorf("Index(sbl) returns %d, should be %d", idx, i)
		}
	}

	nd := New(5)
	_, err := root.Index(nd)
	if err == nil {
		t.Errorf("Index(%q) returns no error, should be %q",
			nd.String(), ErrNodeNotFound.Error())
	} else if !errors.Is(err, ErrNodeNotFound) {
		t.Errorf("Index(nd) returns error %q, should be %q",
			err.Error(), ErrNodeNotFound.Error())
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
				got += fmt.Sprintf("%s%d", space, nd.Get())
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
