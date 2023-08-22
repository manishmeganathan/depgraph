package depgraph

import (
	"testing"

	"github.com/sarvalabs/go-polo"
	"github.com/stretchr/testify/require"
)


func TestDependencyGraph_Serialization_POLO(t *testing.T) {
	inputs := []node{
		{0, nil},
		{1, []uint64{0, 2, 4, 8}},
		{2, nil},
		{3, nil},
		{4, nil},
		{5, nil},
		{6, nil},
		{7, nil},
		{8, []uint64{7, 5, 9}},
		{9, []uint64{0, 4, 5}},
		{10, []uint64{0, 6}},
	}

	dgraph := NewDependencyGraph()
	for _, input := range inputs {
		dgraph.Insert(input.ptr, input.deps...)
	}

	encoded, err := polo.Polorize(dgraph)
	require.Nil(t, err)

	decoded := new(DependencyGraph)
	err = polo.Depolorize(decoded, encoded)
	require.Nil(t, err)

	require.Equal(t, dgraph, decoded)
}
