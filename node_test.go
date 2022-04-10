package otree

import (
	"errors"
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
	tr := New()
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
	tr := New()
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
	tr := New()
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
