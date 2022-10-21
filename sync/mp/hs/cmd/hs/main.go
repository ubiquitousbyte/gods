package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/ubiquitousbyte/gods/sync/mp/hs"
)

func main() {
	nodes := flag.Int("nodes", 32, "Number of nodes in the ring")
	flag.Parse()

	ring := hs.NewRing(*nodes)
	leader := ring.ElectLeader(context.Background())
	fmt.Fprintf(os.Stdout, "Elected leader: %d", leader)
}
