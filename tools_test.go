package otree

import (
	"fmt"
	"testing"
)

func TestInsert(t *testing.T) {
	tests := []struct {
		l1, l2 int
		i      int
		want   string
	}{
		{2, 2, 3, "0 1 2 3"},
		{2, 2, 0, "2 3 0 1"},
		{2, 2, -1, "2 3 0 1"},
		{2, 2, 1, "0 2 3 1"},
		{2, 1, 3, "0 1 2"},
		{2, 1, 0, "2 0 1"},
		{2, 1, 1, "0 2 1"},
		{1, 1, 2, "0 1"},
		{1, 1, 0, "1 0"},
		{0, 1, 1, "0"},
		{0, 1, -1, "0"},
		{1, 0, 2, "0"},
		{1, 0, 1, "0"},
		{1, 0, 0, "0"},
		{0, 0, 0, ""},
		{0, 0, 1, ""},
		{0, 0, -1, ""},
	}

	for _, tst := range tests {
		i1 := make([]*Node, tst.l1)
		i2 := make([]*Node, tst.l2)
		for i := 0; i < tst.l1+tst.l2; i++ {
			if i < tst.l1 {
				i1[i] = NewNode(i)
			} else {
				i2[i-tst.l1] = NewNode(i)
			}
		}
		r := insertNodes(i1, i2, tst.i)
		got := ""
		sep := ""
		for _, nd := range r {
			got += sep + fmt.Sprintf("%d", nd.data.(int))
			sep = " "
		}
		if got != tst.want {
			t.Errorf("InsertSlice() generates %q, should be %q",
				got, tst.want)
		}
	}
}

func TestInvertSlice(t *testing.T) {
	tests := []int{5, 4, 1, 0}

	for _, tst := range tests {
		nodes := []*Node{}
		for i := 0; i < tst; i++ {
			nodes = append(nodes, NewNode(i))
		}

		got := invertSlice(nodes)
		if lGot := len(got); lGot != tst {
			t.Errorf("len(invertSlice(nodes)) returns %d, should be %d, ", lGot, tst)
		} else {
			for i, nd := range got {
				if nd.Get() != tst-i-1 {
					t.Errorf("invertSlice(nodes)[%d] is %d, should be %d",
						i, nd.Get(), tst-i-1)
				}
			}
		}
	}
}

func TestMergePaths(t *testing.T) {
	tests := []struct {
		lA, lB int // lengths of slices to be merged
		lC     int // number of shared nodes
		want   string
		// nodes, b, want []int
	}{
		{5, 5, 2, "[4 3 2 1 5 6 7]"},
		{5, 2, 2, "[4 3 2 1]"},
		{2, 5, 2, "[1 2 3 4]"},
		{5, 1, 1, "[4 3 2 1 0]"},
		{1, 5, 1, "[0 1 2 3 4]"},
		{5, 0, 0, "[]"},
		{0, 5, 0, "[]"},
		{0, 0, 0, "[]"},
	}

	for _, tst := range tests {
		n := 0

		a := make([]*Node, tst.lA)
		for i := 0; i < tst.lA; i++ {
			a[tst.lA-i-1] = NewNode(n)
			n++
		}

		b := make([]*Node, tst.lB)
		c := 0
		for i := 0; i < tst.lB; i++ {
			if c < tst.lC {
				b[tst.lB-i-1] = a[tst.lA-i-1]
				c++
			} else {
				b[tst.lB-i-1] = NewNode(n)
				n++
			}
		}

		got := mergePaths(a, b)

		s := "["
		space := ""
		for _, nd := range got {
			s += fmt.Sprintf("%s%d", space, nd.Get())
			space = " "
		}
		s += "]"
		if s != tst.want {
			t.Errorf("mergPaths(a,b) results in %q, should be %q", s, tst.want)
		}
	}
}
