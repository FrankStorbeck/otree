package otree

import (
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
