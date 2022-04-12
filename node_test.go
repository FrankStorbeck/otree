package otree

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewNode(t *testing.T) {
	s0 := "s0"
	nd := NewNode(s0)
	if got := nd.data.(string); got != s0 {
		t.Errorf("NewNode(%q) returns data %q, should be %q", s0, got, s0)
	}
}

func TestSetGet(t *testing.T) {
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
	nd := NewNode("")
	for _, tst := range tests {
		nd.Set(tst.data)
		got := nd.Get()
		if got != tst.data {
			t.Errorf("Get() after Set(%v) returns %v, should be %v", tst.data, got, tst.data)
		}
	}
}

func TestParent(t *testing.T) {
	tr := New("root")
	s0 := NewNode("s0")
	tr.LinkChildren(tr.root, -1, s0)

	tests := []struct {
		node, parent *Node
		err          error
	}{
		{s0, tr.root, nil},
		{tr.root, nil, ErrParentMissing},
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

func TestIndex(t *testing.T) {
	tr := New("root")
	sbls := []*Node{NewNode(0), NewNode(1), NewNode(2), NewNode(3), NewNode(4)}
	tr.LinkChildren(tr.root, -1, sbls...)
	for i, sbl := range sbls {
		idx, err := tr.root.Index(sbl)
		if err != nil {
			t.Errorf("Index(sbl) returns error %q, should be nil",
				err.Error())
		} else if i != idx {
			t.Errorf("Index(sbl) returns %d, should be %d", idx, i)
		}
	}

	nd := NewNode(5)
	_, err := tr.root.Index(nd)
	if err == nil {
		t.Errorf("Index(sbl) returns no error, should be %q",
			ErrNoNodeFound.Error())
	} else if !errors.Is(err, ErrNoNodeFound) {
		t.Errorf("Index(nd) returns error %q, should be %q",
			err.Error(), ErrNoNodeFound.Error())
	}
}

func TestLevel(t *testing.T) {
	tr := New("root")
	child := NewNode(0)
	grandChild := NewNode(1)

	tr.LinkChildren(tr.root, 0, child)
	tr.LinkChildren(child, 0, grandChild)

	if got := tr.root.Level(); got != 0 {
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
	tr := New("root")
	child := NewNode(0)
	grandChild := NewNode(1)

	tr.LinkChildren(tr.root, 0, child)
	tr.LinkChildren(child, 0, grandChild)

	a := grandChild.Ancestors()
	if got := len(a); got != 2 {
		t.Errorf("len(grandChild.Ancestors()) returns %d, should be 2", got)
	}
	if a[0] != child {
		t.Errorf("grandChild.Ancestors()[0] is not child")
	}
	if a[1] != tr.root {
		t.Errorf("grandChild.Ancestors()[1] is not root")
	}
}

func TestHeight(t *testing.T) {
	tr := New("root")
	child := NewNode(0)
	grandChild := NewNode(1)

	tr.LinkChildren(tr.root, 0, child)
	tr.LinkChildren(child, 0, grandChild)

	if got := tr.root.Height(); got != 2 {
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
	tr := New("root")
	children := []*Node{NewNode(10), NewNode(11), NewNode(12)}
	grandChildren1 := []*Node{NewNode(20), NewNode(21), NewNode(22)}
	greatGrandChildren1 := []*Node{NewNode(30), NewNode(31), NewNode(32)}

	tr.LinkChildren(tr.root, 0, children...)
	tr.LinkChildren(children[0], 0, grandChildren1...)
	tr.LinkChildren(grandChildren1[0], 0, greatGrandChildren1...)

	tests := []struct {
		start, end *Node
		want       string
		distance   int
	}{
		{grandChildren1[2], greatGrandChildren1[0], "[22 10 20 30]", 3},
		{greatGrandChildren1[0], grandChildren1[2], "[30 20 10 22]", 3},
		{greatGrandChildren1[0], children[0], "[30 20 10]", 2},
		{greatGrandChildren1[0], children[1], "[30 20 10 root 11]", 4},
		{children[0], children[0], "[10]", 0},
		{children[0], NewNode(-1), "[]", -1},
		{NewNode(-1), children[0], "[]", -1},
	}

	for _, tst := range tests {
		path := tst.start.Path(tst.end)

		got := "["
		space := ""
		for _, nd := range path {
			switch k := nd.Get().(type) {
			case int:
				got += fmt.Sprintf("%s%d", space, k)
			case string:
				got += fmt.Sprintf("%s%s", space, k)
			}
			space = " "
		}
		got += "]"

		if got != tst.want {
			t.Errorf("Path() results in %q, should be %q", got, tst.want)
		} else if d := tst.start.Distance(tst.end); d != tst.distance {
			t.Errorf("Distance() is %d, should be %d", d, tst.distance)
		}
	}
}
