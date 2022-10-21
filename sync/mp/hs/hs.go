// Hirschberg and Sinclair (HS) algorithm
//
// Reduces the complexity of lcr from O(n) down to O(n log n)
//
// Does not rely on n, supports bidirectional communication.
// Relies on the fact that node identifiers are consecutive integers.
//
// Each process operates in phases 0, 1, 2, ... In each phase l, process i
// sends tokens containing its identifier in both directions.
// These travel a certain distance, 2^l, and return back to their origin.
// If both tokens make it back safely, i continues with the following phase.
// However, tokens might not macke it back safely. While a u_i token
// is proceeding in the outbound direction, each process j on u_i's path
// compares u_i with its own identifier u_j. If u_i < u_j, j discards the token,
// whereas if u_i > u_j, j relays u_i. If u_i = u_j, then process j has
// received its own identifier before the token has turned around, so it
// elects itself the leader
package hs

import (
	"context"
	"math"
)

type Direction int

const (
	DirectionOut Direction = iota
	DirectionIn
)

type Status int

const (
	StatusUnknown Status = iota
	StatusLeader
)

type Ring []*Node

func NewRing(count int) Ring {
	if count <= 0 {
		return nil
	}

	first := newNode(0, nil, nil)
	r := []*Node{first}
	for i := 1; i < count; i++ {
		prev := r[(i-1)%count]
		r = append(r, newNode(i, prev.txLeft, prev.txRight))
	}
	first.rxLeft, first.rxRight = r[count-1].txLeft, r[count-1].txRight

	return r
}

func (r Ring) ElectLeader(ctx context.Context) int {
	leader := make(chan int)
	defer close(leader)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, node := range r {
		go func(ctx context.Context, n *Node) {
			defer close(n.txLeft)
			defer close(n.txRight)
			for {
				select {
				case <-ctx.Done():
					return
				default:
					if n.IsLeader() {
						leader <- n.id
						return
					}
					n.RoundTrip()
				}
			}
		}(ctx, node)
	}

	return <-leader
}

type Node struct {
	id                  int
	sendLeft, sendRight *Token
	txLeft, txRight     chan Token
	rxLeft, rxRight     chan Token
	phase               int
	status              Status
}

func newNode(id int, rxLeft, rxRight chan Token) *Node {
	return &Node{
		id:        id,
		sendLeft:  newToken(id),
		sendRight: newToken(id),
		txLeft:    make(chan Token, 1),
		txRight:   make(chan Token, 1),
		rxLeft:    rxLeft,
		rxRight:   rxRight,
	}
}

func (n *Node) RoundTrip() {
	n.generate()
	n.transition()
}

func (n *Node) IsLeader() bool {
	return n.status == StatusLeader
}

func (n *Node) generate() {
	if tok := n.sendLeft; tok != nil {
		n.txLeft <- *tok
	}

	if tok := n.sendRight; tok != nil {
		n.txRight <- *tok
	}
}

func (n *Node) transition() {
	n.sendLeft, n.sendRight = nil, nil

	// Receive a token from the left neighbour
	tokLeft, ok := <-n.rxLeft
	if !ok {
		return
	}

	if tokLeft.direction == DirectionOut {
		if tokLeft.id > n.id {
			if tokLeft.hops > 1 {
				// Keep propagating the token to the right
				n.sendRight = &tokLeft
				n.sendRight.hops -= 1
			} else if tokLeft.hops == 1 {
				// Reverse the direction of the token
				n.sendLeft = &Token{id: tokLeft.id, direction: DirectionIn, hops: 1}
			}
		} else if tokLeft.id == n.id {
			n.status = StatusLeader
		}
	} else {
		// Token is being relayed back
		if tokLeft.id != n.id && tokLeft.hops == 1 {
			n.sendRight = &tokLeft
		}
	}

	// Receive a token from the right neighbour
	tokRight, ok := <-n.rxRight
	if !ok {
		return
	}

	if tokRight.direction == DirectionOut {
		if tokRight.id > n.id {
			if tokRight.hops > 1 {
				// The token must be propagated to the left until its direction
				// gets reversed
				n.sendLeft = &tokRight
				n.sendLeft.hops -= 1
			} else if tokRight.hops == 1 {
				// Time to reverse the token's direction
				n.sendRight = &Token{id: tokRight.id, direction: DirectionIn, hops: 1}
			}
		} else if tokRight.id == n.id {
			n.status = StatusLeader
		}
	} else {
		// Token is being relayed back
		if tokRight.id != n.id && tokRight.hops == 1 {
			n.sendLeft = &tokRight
		}
	}

	if tokLeft == tokRight && tokLeft.id == n.id && tokLeft.direction == DirectionIn {
		n.phase += 1
		hops := int(math.Pow(2, float64(n.phase)))
		n.sendRight = &Token{id: n.id, hops: hops}
		n.sendLeft = &Token{id: n.id, hops: hops}
	}
}

type Token struct {
	id        int
	hops      int
	direction Direction
}

func newToken(id int) *Token {
	return &Token{id: id, hops: 1}
}
