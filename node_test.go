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
