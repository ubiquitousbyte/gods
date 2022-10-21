# Synchronous message-passing algorithms

Some are more complex than others. In particular, certain algorithms
are restricted to work on ring topologies, e.g the 
Le Lann, Chang and Roberts algorithm, as well as the Hirschberg and Sinclair
algorithm. 

## Rings

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

## General networks
G is a strongly connected network digraph with n nodes. Nodes are 
named [1, n] but algorithms are not allowed to rely on their 
identifiers being consecutive integers.

### Additional challenges
1. Same as 1. in Rings
2. The number of nodes and the diamater 
(the maximum distance `dist(i, j)` taken over all process pairs `(i, j)`) may 
can be unknown. 