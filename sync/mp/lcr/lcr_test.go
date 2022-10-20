package lcr

import (
	"context"
	"testing"
)

func TestElectLeader(t *testing.T) {
	for i := 0; i < 16; i++ {
		nodes := 1 << i
		ring := NewRing(nodes)
		got := ring.ElectLeader(context.TODO())

		exp := nodes - 1
		if got != exp {
			t.Fatalf("expected leader to be %d, but got %d", exp, got)
		}

		for _, node := range ring[:exp] {
			if node.status != StatusNonLeader {
				t.Fatalf("expected non leader status for node %d", node.id)
			}
		}
	}
}
