package main

import (
	"reflect"
	"testing"
)

// A simple implemention for Vertex type for testing
type IntVertex int

func (v IntVertex) ID() VertexID {
	return int(v)
}

func TestAddVertex(t *testing.T) {
	for _, testCase := range []struct {
		vertices      []IntVertex
		expectedCount int
	}{
		{
			vertices:      []IntVertex{},
			expectedCount: 0,
		},
		{
			vertices:      []IntVertex{IntVertex(1)},
			expectedCount: 1,
		},
		{
			vertices:      []IntVertex{IntVertex(1), IntVertex(1)},
			expectedCount: 1,
		},
		{
			vertices:      []IntVertex{IntVertex(1), IntVertex(2)},
			expectedCount: 2,
		},
	} {
		g := NewGraph()
		for _, v := range testCase.vertices {
			g.Add(v)
		}
		if len(g.vertices) != testCase.expectedCount {
			t.Errorf("expected count: %d, actual: %v", testCase.expectedCount, g.vertices)
		}
	}
}

func TestLinkBoth(t *testing.T) {
	g := NewGraph()
	g.LinkBoth(IntVertex(1), IntVertex(2), 1)
	g.LinkBoth(IntVertex(2), IntVertex(3), 2)

	expected := map[VertexID]map[VertexID]EdgeWeight{
		1: map[VertexID]EdgeWeight{
			2: EdgeWeight(1),
		},
		2: map[VertexID]EdgeWeight{
			1: EdgeWeight(1),
			3: EdgeWeight(2),
		},
		3: map[VertexID]EdgeWeight{
			2: EdgeWeight(2),
		},
	}

	if !reflect.DeepEqual(expected, g.edges) {
		t.Errorf("expected: %v, actual: %v", expected, g.edges)
	}
}
