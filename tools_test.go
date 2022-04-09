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
