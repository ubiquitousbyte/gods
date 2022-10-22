// Every process keeps a record of the maximum identifier it has seen so far.
// At each round, each process "floods" this maximum to all of its peers.
// After diam rounds, if the maximum value recorded by the process is its
// own identifier, it elects itself the leader
package floodmax

import (
	"sync"
)

type Status int

const (
	StatusUnknown Status = iota
	StatusNonLeader
	StatusLeader
)

type Network []*Node

type Node struct {
	id, max, rounds int
	status          Status
	diam            int
	txs, rxs        []chan int
}

func newNode(id int) *Node {
	return &Node{id: id, max: id, txs: make([]chan int, 3)}
}

func (n *Node) IsLeader() bool {
	return n.status == StatusLeader
}

func (n *Node) RoundTrip() {
	var wg sync.WaitGroup

	wg.Add(len(n.txs))
	for _, tx := range n.txs {
		// Flood the network with the maximum value
		go func(tx chan<- int) {
			tx <- n.max
			wg.Done()
		}(tx)
	}
	wg.Wait()

	n.rounds += 1

	maxes := make(chan int, len(n.rxs))

	wg.Add(len(n.rxs))
	// Fan-in the channels
	for _, rx := range n.rxs {
		go func(rx <-chan int) {
			defer wg.Done()
			n, ok := <-rx
			if !ok {
				return
			}
			maxes <- n
		}(rx)
	}
	wg.Wait()

	close(maxes)

	for m := range maxes {
		if m > n.max {
			n.max = m
		}
	}

	if n.rounds == n.diam {
		if n.max == n.id {
			n.status = StatusLeader
		} else {
			n.status = StatusNonLeader
		}
	}
}
