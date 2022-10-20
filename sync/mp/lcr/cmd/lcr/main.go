package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/ubiquitousbyte/gods/sync/mp/lcr"
)

func main() {
	nodes := flag.Int("nodes", 32, "Number of nodes in the network")
	flag.Parse()

	ring := lcr.NewRing(*nodes)
	leader := ring.ElectLeader(context.Background())
	fmt.Fprintf(os.Stdout, "Elected leader: %d\n", leader+1)
}
