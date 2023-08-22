package depgraph

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	mapset "github.com/deckarep/golang-set"
)

type DependencyGraph struct {
	mutex sync.RWMutex
	graph map[uint64]mapset.Set
}

// NewDependencyGraph generates and returns an empty DependencyGraph
func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{graph: make(map[uint64]mapset.Set)}
}

// Insert inserts an uint64 as a graph vertex to the DependencyGraph.
// It also accepts a variadic number of dependencies for the pointer and inserts them as edges.
//
// If the vertex (and subsequently its edges) already exists, it is overwritten.
func (dgraph *DependencyGraph) Insert(ptr uint64, deps ...uint64) {
	// Create a new Set and insert the dependencies into it
	set := mapset.NewSet()
	for _, dep := range deps {
		set.Add(dep)
	}

	// Push the vertex and its edges into the graph
	dgraph.push(ptr, set)
}

// Remove removes an uint64 as a graph vertex from the DependencyGraph.
// If such a vertex does not exist, this is a no-op.
func (dgraph *DependencyGraph) Remove(ptr uint64) {
	dgraph.pop(ptr)
}

// Size returns the number of vertices in the DependencyGraph
func (dgraph *DependencyGraph) Size() uint64 {
	dgraph.mutex.RLock()
	defer dgraph.mutex.RUnlock()

	return uint64(len(dgraph.graph))
}

// Contains returns whether a given vertex exists in the DependencyGraph
func (dgraph *DependencyGraph) Contains(ptr uint64) bool {
	_, ok := dgraph.peek(ptr)

	return ok
}

// Copy creates a clone of the DependencyGraph and returns it
func (dgraph *DependencyGraph) Copy() *DependencyGraph {
	dgraph.mutex.RLock()
	defer dgraph.mutex.RUnlock()

	// Create a DependencyGraph with a graph buffer large enough for all elements in the original
	clone := &DependencyGraph{graph: make(map[uint64]mapset.Set, len(dgraph.graph))}

	// For each vertex pointer, copy its edge dependencies and insert
	for ptr, deps := range dgraph.graph {
		clone.push(ptr, deps.Clone())
	}

	return clone
}

// String implements the Stringer interface for DependencyGraph.
func (dgraph *DependencyGraph) String() string {
	dgraph.mutex.RLock()
	defer dgraph.mutex.RUnlock()

	elements := make([]string, 0)

	// Iterate over the graph vertices
	for ptr := range dgraph.Iter() {
		// Get the edge dependencies
		deps, _ := dgraph.peek(ptr)
		// If no edges, just add the pointer value
		if deps.Cardinality() == 0 {
			elements = append(elements, fmt.Sprintf("%v", ptr))

			continue
		}

		// Sort the deps and format element as ptr:[deps]
		depSlice := deps.ToSlice()
		sort.Slice(depSlice, func(i, j int) bool {
			return depSlice[i].(uint64) < depSlice[j].(uint64) //nolint:forcetypeassert
		})

		elements = append(elements, fmt.Sprintf("%v:%v", ptr, depSlice))
	}

	return fmt.Sprintf("DependencyGraph{%v}", strings.Join(elements, ", "))
}

// Iter returns a channel iterator that iterates over the vertices of the DependencyGraph is sorted order.
// This iteration is thread-safe, the graph being immutable during the iteration.
func (dgraph *DependencyGraph) Iter() <-chan uint64 {
	ch := make(chan uint64)

	go func() {
		dgraph.mutex.RLock()
		defer dgraph.mutex.RUnlock()

		for _, ptr := range dgraph.sorted() {
			ch <- ptr
		}

		close(ch)
	}()

	return ch
}

// Edges returns the edges of going out of a given vertex pointer.
// The dependencies are returned as a mapset.Set (cardinality is zero if no dependencies for vertex)
func (dgraph *DependencyGraph) Edges(ptr uint64) []uint64 {
	depSet, _ := dgraph.peek(ptr)

	deps := make([]uint64, 0, depSet.Cardinality())
	for dep := range depSet.Iter() {
		deps = append(deps, dep.(uint64)) //nolint:forcetypeassert
	}

	return deps
}

// Dependencies returns all the edges (and edges of edges) for a given vertex pointer.
// It recursively collects all dependencies from each dependency layer and returns them (without duplicates).
// Note: This should only be used if the DependencyGraph can be resolved, otherwise, it will result in an infinite loop.
func (dgraph *DependencyGraph) Dependencies(ptr uint64) []uint64 {
	depSet := mapset.NewSet()

	// Collect all the direct deps of the pointer
	for _, dep := range dgraph.Edges(ptr) {
		// Add the direct dep to depSet
		depSet.Add(dep)

		// Recursively collect all sub dependencies
		deeper := dgraph.Dependencies(dep)
		if len(deeper) == 0 {
			continue
		}

		// Add all sub dependencies to the set
		for _, dep := range deeper {
			depSet.Add(dep)
		}
	}

	// Collect all dependencies (free from duplicates)
	deps := make([]uint64, 0, depSet.Cardinality())
	for dep := range depSet.Iter() {
		deps = append(deps, dep.(uint64)) //nolint:forcetypeassert
	}

	// Sort the dependencies
	sort.Slice(deps, func(i, j int) bool {
		return deps[i] < deps[j]
	})

	return deps
}
