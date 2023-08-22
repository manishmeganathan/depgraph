package depgraph

import (
	"sort"

	mapset "github.com/deckarep/golang-set"
)

// ResolveBatches attempts to resolve the DependencyGraph into batched element pointers.
// Each batch represents elements that need to compiled before the next batch but are independent of
// each other. The output of graph resolution is deterministic as each batch of pointers is sorted.
//
// Returns a boolean along with the batches indicating if the graph could be resolved.
// Graph resolution fails if there are circular or nil (non-existent) dependencies.
func (dgraph *DependencyGraph) ResolveBatches() ([][]uint64, bool) {
	dgraph.mutex.RLock()
	defer dgraph.mutex.RUnlock()

	// Create a working copy of the graph
	working := dgraph.Copy()
	// Initialize the slice of element batches
	batches := make([][]uint64, 0)

	// Iterate until, the working graph has been emptied
	for working.Size() != 0 {
		ready := mapset.NewSet()

		// Accumulate all elements from the working
		// set that have zero unresolved dependencies
		for ptr := range working.Iter() {
			if deps, _ := working.peek(ptr); deps.Cardinality() == 0 {
				ready.Add(ptr)
			}
		}

		// If there are no ready elements, we have an issue
		// Either a circular or nil dependency exists in the graph
		if ready.Cardinality() == 0 {
			return nil, false
		}

		// Remove all the elements that are ready from the working graph
		for item := range ready.Iter() {
			working.pop(item.(uint64)) //nolint:forcetypeassert
		}

		// Remove the dependencies for each element in the working graph that have now been resolved.
		// We calculate the difference for each remaining set of dependency edges compared to the ready set.
		for ptr := range working.Iter() {
			deps, _ := working.peek(ptr)
			working.graph[ptr] = deps.Difference(ready)
		}

		// Accumulate all pointers in the ready set as an element batch
		batch := make([]uint64, 0, ready.Cardinality())
		for item := range ready.Iter() {
			batch = append(batch, item.(uint64)) //nolint:forcetypeassert
		}

		// Sort the element batch (the order within a batch is not
		// important, but we do this to get a deterministic output)
		sort.Slice(batch, func(i, j int) bool {
			return batch[i] < batch[j]
		})

		batches = append(batches, batch)
	}

	return batches, true
}

// Resolve attempts to resolve the DependencyGraph into an ordered slice of element pointers.
// This slice represents the order of element compilation and is always deterministic.
// The output our Resolve is essentially a flattened output of ResolveBatches.
//
// Returns a boolean along with the resolved elements indicating if the graph could be resolved.
// Graph resolution fails if there are circular or nil (non-existent) dependencies.
func (dgraph *DependencyGraph) Resolve() ([]uint64, bool) {
	dgraph.mutex.RLock()
	defer dgraph.mutex.RUnlock()

	// Resolve the graph into batched elements
	batches, ok := dgraph.ResolveBatches()
	if !ok {
		return nil, false
	}

	// Flatten the batches into a single slice
	// The output inherits its determinism from ResolveBatches
	resolved := make([]uint64, 0, dgraph.Size())
	for _, batch := range batches {
		resolved = append(resolved, batch...)
	}

	return resolved, true
}
