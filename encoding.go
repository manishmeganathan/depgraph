package depgraph

import (
	"encoding/json"
	"sort"

	"github.com/sarvalabs/go-polo"
)

// MarshalJSON implements the json.Marshaller interface for DependencyGraph
func (dgraph *DependencyGraph) MarshalJSON() ([]byte, error) {
	// Get encodable version of dgraph
	encodable := dgraph.encode()

	return json.Marshal(encodable)
}

// UnmarshalJSON implements the json.Unmarshaller interface for DependencyGraph
func (dgraph *DependencyGraph) UnmarshalJSON(data []byte) error {
	// Decode the data into a map of graph nodes
	decodable := make(map[uint64][]uint64)
	if err := json.Unmarshal(data, &decodable); err != nil {
		return err
	}

	// Decode data into dgraph
	dgraph.decode(decodable)

	return nil
}

// Polorize implements the polo.Polorizable interface for DependencyGraph
func (dgraph *DependencyGraph) Polorize() (*polo.Polorizer, error) {
	// Get encodable version of dgraph
	encodable := dgraph.encode()

	// Serialize the encodable map
	polorizer := polo.NewPolorizer()
	if err := polorizer.Polorize(encodable); err != nil {
		return nil, err
	}

	return polorizer, nil
}

// Depolorize implements the polo.Depolorizable interface for DependencyGraph
func (dgraph *DependencyGraph) Depolorize(depolorizer *polo.Depolorizer) error {
	// Decode the data into a map of graph nodes
	decodable := make(map[uint64][]uint64)
	if err := depolorizer.Depolorize(&decodable); err != nil {
		return err
	}

	// Decode data into dgraph
	dgraph.decode(decodable)

	return nil
}

// encode converts the DependencyGraph into a map of pointers to their dependencies.
// The generated map is safe to encode with any encoding scheme.
func (dgraph *DependencyGraph) encode() map[uint64][]uint64 {
	// Declare a map to collect all graph nodes
	encodable := make(map[uint64][]uint64, dgraph.Size())
	// Iterate over the graph vertices
	for ptr := range dgraph.Iter() {
		// Get the edge dependencies
		depSet, _ := dgraph.peek(ptr)

		// Collect the edges into a []uint64
		deps := make([]uint64, 0, depSet.Cardinality())
		for dep := range depSet.Iter() {
			deps = append(deps, dep.(uint64)) //nolint:forcetypeassert
		}

		// Sort the dependencies
		sort.Slice(deps, func(i, j int) bool {
			return deps[i] < deps[j]
		})

		encodable[ptr] = deps
	}

	return encodable
}

// decode converts a given map of pointers to their dependencies into a DependencyGraph and absorbs it.
func (dgraph *DependencyGraph) decode(data map[uint64][]uint64) {
	// Insert each node into the graph
	*dgraph = *NewDependencyGraph()

	for ptr, deps := range data {
		dgraph.Insert(ptr, deps...)
	}
}
