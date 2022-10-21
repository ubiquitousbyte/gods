// Hirschberg and Sinclair (HS) algorithm
//
// Reduces the complexity of lcr from O(n) down to O(n log n)
//
// Does not rely on n, supports bidirectional communication.
// Relies on the fact that node identifiers are consecutive integers.
package hs

import "context"

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

type Token struct {
	id        int
	hops      int
	direction Direction
}

func newToken(id int) *Token {
	return &Token{id: id, hops: 1}
}

type Node struct {
	id           int
	sendp, sendm *Token
	txp, rxp     chan Token
	txm, rxm     chan Token
	status       Status
	phase        int
}

func newNode(id int, rxp, rxm chan Token) *Node {
	return &Node{
		id:    id,
		sendp: newToken(id),
		sendm: newToken(id),
		txp:   make(chan Token, 1),
		txm:   make(chan Token, 1),
		rxp:   rxp,
		rxm:   rxm,
	}
}

func (n *Node) generate() {
	if tok := n.sendp; tok != nil {
		n.txp <- *tok
	}

	if tok := n.sendm; tok != nil {
		n.txm <- *tok
	}
}

func (n *Node) transition() {

}

func (n *Node) RoundTrip(ctx context.Context) {
	n.generate()
	n.transition()
}
