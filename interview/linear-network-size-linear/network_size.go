// Code and test for blog
// https://peng.fyi/post/interview/calculate-linear-network-size-by-message-passing/
package networksize

type Node struct {
	Left  *Node
	Right *Node

	// nodes on the left side and right side.
	nodesLeft  int
	nodesRight int
}

func (node *Node) sendLeft(m int) {
	if node.Left != nil {
		node.Left.receiveRight(m)
	}
}

func (node *Node) sendRight(m int) {
	if node.Right != nil {
		node.Right.receiveLeft(m)
	}
}

// L -> A. node A receives information from node L.
// if m > 0: L tells A that there are m nodes on A's left side, including L.
// otherwise: L asks A for how many nodes are on L's right side, including A.
func (node *Node) receiveLeft(m int) {
	if m > 0 {
		node.nodesLeft = m
		node.sendRight(node.nodesLeft + 1)
		return
	}

	if node.Right == nil {
		node.sendLeft(1)
	} else if node.nodesRight > 0 {
		node.sendLeft(node.nodesRight + 1)
	} else {
		node.sendRight(0)
	}
}

// A <- R. node A receives information from node R.
// if m > 0: R tells A that there are n nodes on A's right side, including R.
// otherwise: R asks A for how many nodes are on R's left side, including A.
func (node *Node) receiveRight(m int) {
	if m > 0 {
		node.nodesRight = m
		node.sendLeft(node.nodesRight + 1)
		return
	}

	if node.Left == nil {
		node.sendRight(1)
	} else if node.nodesLeft > 0 {
		node.sendRight(node.nodesLeft + 1)
	} else {
		node.sendLeft(0)
	}
}

func (node *Node) GetNetworkSize() int {
	if node.Left == nil && node.Right == nil {
		return 1
	}

	if node.Left == nil {
		node.sendRight(1)
		node.sendRight(0)
	}
	if node.Right == nil {
		node.sendLeft(1)
		node.sendLeft(0)
	}
	if node.Left != nil && node.Right != nil {
		node.sendLeft(0)
		node.sendRight(0)
	}

	return node.nodesLeft + node.nodesRight + 1
}
