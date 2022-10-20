# Distributed Systems in Go

Notes and code that articulate certain problems in distributed computing.

Algorithms are classified based on two properties - timing and inter-process communication. The package structure aims 
to mirror the aforementioned classification. Each root-level package depicts the timing properties of the algorithms that reside within, e.g `sync` contains only synchronous algorithms. Within each root-level package, you'll find subpackages denoting the inter-process communication model. For example, package `sync/shm`
contains synchronous shared-memory algorithms, whereas `async/shmsg` contains asynchronous message-passing algorithms. 