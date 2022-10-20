// LCR (Le Lann, Chang and Roberts) algorithm
//
// Does not rely on n, supports only unidirectional communication.
// Non-leader nodes output that they are not leading.
// Relies on the fact that node identifiers are consecutive integers.
//
// Each node sends its identifier around the ring. When a process receives
// an incoming identifier, it compares it to its own. If the incoming
// identifier is greater, it keeps propagating it to its neighbour.
// If it is smaller, it discards it. If they are equal, then the node
// denotes itself the leader.
package lcr

import (
	"context"
)

type Status int

const (
	StatusUnknown Status = iota
	StatusNonLeader
	StatusLeader
)

type Ring []*Node

func NewRing(count int) Ring {
	if count <= 0 {
		return nil
	}

	first := newNode(0, nil)
	r := []*Node{first}
	for i := 1; i < count; i++ {
		r = append(r, newNode(i, r[(i-1)%count].tx))
	}
	first.rx = r[count-1].tx

	return r
}

func (r Ring) ElectLeader(ctx context.Context) int {
	leader := make(chan int)
	defer close(leader)

	ctx, cancel := context.WithCancel(ctx)
	// Note that Lynch suggests to implement halting by having the leader
	// send out a special report message after being elected.
	// Any process that receives the report message can halt, after
	// passing the message on to its neighbour.
	//
	// We use context cancelation to avoid the verbosity.
	// The main goroutine acts as a halter that notifies all nodes
	// by canceling the context. The nodes, in turn, monitor the context's
	// done channel after every round and exit when the channel is closed.
	defer cancel()

	for _, node := range r {
		go func(ctx context.Context, n *Node) {
			defer close(n.tx)
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
	id     int
	sendId *int
	status Status
	rx     chan int
	tx     chan int
}

func newNode(id int, rx chan int) *Node {
	return &Node{
		id:     id,
		sendId: &id,
		tx:     make(chan int, 1),
		rx:     rx,
	}
}

func (n *Node) IsLeader() bool {
	return n.status == StatusLeader
}

func (n *Node) RoundTrip() {
	if id := n.sendId; id != nil {
		n.tx <- *id
	}

	id, ok := <-n.rx
	if !ok {
		return
	}

	n.sendId = nil

	if id == n.id {
		n.status = StatusLeader
	} else if id > n.id {
		// If a node detects a larger identifier,
		// then it knows it isn't the leader
		n.status = StatusNonLeader
		n.sendId = &id
	}
}
