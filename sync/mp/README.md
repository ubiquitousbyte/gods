# Synchronous message-passing algorithms

Some are more complex than others. In particular, certain algorithms
are restricted to work on ring topologies, e.g the 
Le Lann, Chang and Roberts algorithm, as well as the Hirschberg and Sinclair
algorithm. 

## Ring algorithms

G consists of n nodes [1, n] organised in a ring. 
Nodes in G do not know their indices.
Message generation and state-transition functions have internal naming
schemes to reference neighbours.
Requirement: Eventually, exactly one process outputs that it
is the leader, e.g by changing its state.

### Additional challenges:
1. Non-leader processes output that they are not the leader
2. Ring directionality - If unidirectional, each node has only one
outgoing neighbour, i.e its clockwise neighbour
3. n is not known to the processes, i.e nodes must work in
dynamically-sized rings
4. Processes should not rely on the fact that their identifiers are
consecutive integers
