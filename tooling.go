package depgraph

import (
	"sort"

	mapset "github.com/deckarep/golang-set"
)

// push accepts a vertex as an uint64 and its dependency
// edges as a mapset.Set and inserts them into the DependencyGraph
func (dgraph *DependencyGraph) push(ptr uint64, set mapset.Set) {
	dgraph.mutex.Lock()
	defer dgraph.mutex.Unlock()

	dgraph.graph[ptr] = set
}

// pop accepts a vertex as an uint64 and removes it from the graph
func (dgraph *DependencyGraph) pop(ptr uint64) {
	dgraph.mutex.Lock()
	defer dgraph.mutex.Unlock()

	delete(dgraph.graph, ptr)
}

// peek returns the edge dependencies for a given vertex as a
// mapset.Set along with a boolean indicating if the vertex existed.
func (dgraph *DependencyGraph) peek(ptr uint64) (mapset.Set, bool) {
	dgraph.mutex.RLock()
	defer dgraph.mutex.RUnlock()

	set, ok := dgraph.graph[ptr]

	return set, ok
}

// sorted returns the vertices of the DependencyGraph as a sorted slice of uint64
func (dgraph *DependencyGraph) sorted() []uint64 {
	dgraph.mutex.RLock()
	defer dgraph.mutex.RUnlock()

	sorted := make([]uint64, 0, len(dgraph.graph))
	for ptr := range dgraph.graph {
		sorted = append(sorted, ptr)
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	return sorted
}
