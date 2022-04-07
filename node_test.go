package otree

import "testing"

func TestNewNode(t *testing.T) {
	s0 := "s0"
	nd := NewNode(s0)
	if got := nd.data.(string); got != s0 {
		t.Errorf("NewNode(%q) returs data %q, should be %q", s0, got, s0)
	}
}
