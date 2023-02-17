package otree

import (
	"testing"
)

func TestTreeWidthAndBreadthAndSize(t *testing.T) {
	root := New("root")
	children := []*Node{New(10), New(11), New(12)}
	grandChildren1 := []*Node{New(20), New(21), New(22), New(23)}
	greatGrandChildren1 := []*Node{New(30), New(31)}

	root.Link(0, children...)
	children[1].Link(0, grandChildren1...)
	grandChildren1[2].Link(0, greatGrandChildren1...)

	want := []int{1, 3, 4, 2, 0}
	for i := 0; i < len(want); i++ {
		if w := Width(root, i, false); w != want[i] {
			t.Errorf("Width(%d) returns %d, should be %d", i, w, want[i])
		}
	}

	wantB := 0
	for _, w := range want {
		if w == 0 {
			wantB++
		} else {
			wantB += (w - 1)
		}
	}
	if b := Breadth(root, false); b != wantB {
		t.Errorf("Breadth() returns %d, should be %d", b, wantB)
	}

	wantS := 0
	for _, w := range want {
		wantS += w
	}
	if s := Size(root, false); s != wantS {
		t.Errorf("Size() returns %d, should be %d", s, wantS)
	}
}

func TestTreeDegree(t *testing.T) {
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
		if d := Degree(tst.root, false); d != tst.degree {
			t.Errorf("%d: Degree() returns %d, should be %d",
				i, d, tst.degree)
		}
	}
}

func TestTreeHeight(t *testing.T) {
	root := New("root")
	child := New(0)
	grandChild := New(1)

	root.Link(0, child)
	child.Link(0, grandChild)

	if got := Height(root.Root(), false); got != 2 {
		t.Errorf("Height() returns %d, should be 2", got)
	}
	if got := Height(child.Root(), false); got != 2 {
		t.Errorf("Height() returns %d, should be 2", got)
	}
	if got := Height(grandChild.Root(), false); got != 2 {
		t.Errorf("Height() returns %d, should be 2", got)
	}
}

func TestRemoveNodeAndRemoveSiblings(t *testing.T) {
	root := New("root")
	children := []*Node{New(10), New(11), New(12)}
	grandChildren := []*Node{New(20), New(21), New(22), New(23)}
	greatGrandChildren := []*Node{New(30), New(31)}

	root.Link(0, children...)
	children[1].Link(0, grandChildren...)
	grandChildren[2].Link(0, greatGrandChildren...)

	tests := []struct {
		node  *Node
		index int
		err   error
		want  string
	}{
		{children[1], AtStart, nil, "<root>[<10>,<12>]"},
		{children[1], AtStart, ErrNodeNotFound, ""},
		{root, AtStart, ErrCannotRemoveRootNode, ""},
		{New("-1"), AtStart, ErrNodeNotFound, ""},
		{nil, AtStart, ErrNodeNotFound, ""},
		{nil, 1, nil, "<root>[<10>]"},
		{nil, AtStart, ErrNodeNotFound, ""},
		{nil, 1, ErrNodeNotFound, ""},
	}

	for _, tst := range tests {
		var err error
		if tst.index < 0 {
			err = Remove(root, tst.node)
		} else {
			tst.node, err = root.RemoveSibling(tst.index)
		}
		switch {
		case err != nil && tst.err == nil:
			t.Errorf("RemoveNode(%v) returns an error %q, should be nil",
				tst.node.Data, err.Error())

		case err == nil && tst.err != nil:
			t.Errorf("RemoveNode(%v) returns no error, should be %q",
				tst.node.Data, tst.err.Error())

		case err != nil && tst.err != nil && err != tst.err:
			t.Errorf("RemoveNode(%v) returns error %q, should be %q",
				tst.node.Data, err.Error(), tst.err.Error())

		case err == nil && tst.err == nil:
			if got := root.String(); got != tst.want {
				t.Errorf("RemoveNode(%v) generates %q, should be %q",
					tst.node.Data, got, tst.want)
			}
		}
	}
}

// func TestMain(t *testing.T) {
// 	root := New("root")
//
// 	level1 := []*Node{New(10), New(11), New(12), New(13)}
// 	root.Link(0, level1...)
//
// 	level2 := []*Node{New(20), New(21)}
// 	level1[0].Link(0, level2...)
//
// 	fmt.Println(root.String()) // output: <root>[<10>[<20>,<21>],<11>,<12>,<13>]
// }
