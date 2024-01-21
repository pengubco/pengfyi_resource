package networksize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinear_1(t *testing.T) {
	cases := []struct {
		nodeCount int
	}{
		{1},
		{2},
		{3},
		{4},
		{5},
		{2000},
	}

	for _, tc := range cases {
		t.Run("", func(tt *testing.T) {
			nodes := createNetworkOfSize(tc.nodeCount)
			for i := 0; i < tc.nodeCount; i++ {
				assert.Equal(tt, tc.nodeCount, nodes[i].GetNetworkSize())
			}
		})
	}
}

// Create a linear network of n servers. Return servers from left to right.
func createNetworkOfSize(n int) []Node {
	nodes := make([]Node, n)
	for i := 0; i < n; i++ {
		if i < n-1 {
			nodes[i].Right = &nodes[i+1]
		}
		if i > 0 {
			nodes[i].Left = &nodes[i-1]
		}
	}
	return nodes
}
