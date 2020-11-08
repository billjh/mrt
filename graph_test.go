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
		if len(g.Vertices) != testCase.expectedCount {
			t.Errorf("expected count: %d, actual: %v", testCase.expectedCount, g.Vertices)
		}
	}
}

func TestLinkBoth(t *testing.T) {
	g := NewGraph()
	g.LinkBoth(IntVertex(1), IntVertex(2), 1)
	g.LinkBoth(IntVertex(2), IntVertex(3), 2)

	expected := map[VertexID]map[VertexID]Weight{
		1: map[VertexID]Weight{
			2: Weight(1),
		},
		2: map[VertexID]Weight{
			1: Weight(1),
			3: Weight(2),
		},
		3: map[VertexID]Weight{
			2: Weight(2),
		},
	}

	if !reflect.DeepEqual(expected, g.Edges) {
		t.Errorf("expected: %v, actual: %v", expected, g.Edges)
	}
}

// test cases that expect error for searching
var expectedError = []struct {
	name string
	g    *Graph
	src  VertexID
	dest VertexID
}{
	{
		name: "empty graph",
		g:    NewGraph(),
		src:  1,
		dest: 2,
	},
	{
		name: "unknown src",
		g:    NewGraph().Add(IntVertex(2)),
		src:  1,
		dest: 2,
	},
	{
		name: "unknown dest",
		g:    NewGraph().Add(IntVertex(1)),
		src:  1,
		dest: 2,
	},
	{
		name: "src and destination are same",
		g:    NewGraph().Add(IntVertex(1)),
		src:  1,
		dest: 1,
	},
	{
		name: "two stand alone vertices",
		g:    NewGraph().Add(IntVertex(1)).Add(IntVertex(2)),
		src:  1,
		dest: 2,
	},
	{
		name: "disconnected graph",
		g:    NewGraph().LinkBoth(IntVertex(1), IntVertex(2), 1).LinkBoth(IntVertex(3), IntVertex(4), 1),
		src:  1,
		dest: 4,
	},
}

func TestBFSError(t *testing.T) {
	for _, testCase := range expectedError {
		_, err := testCase.g.BFS(testCase.src, testCase.dest)
		if err == nil {
			t.Errorf("expect error on %s", testCase.name)
		}
	}
}

func TestDijkstraError(t *testing.T) {
	for _, testCase := range expectedError {
		_, err := testCase.g.Dijkstra(testCase.src, testCase.dest)
		if err == nil {
			t.Errorf("expect error on %s", testCase.name)
		}
	}
}

func TestDijkstraAllError(t *testing.T) {
	for _, testCase := range expectedError {
		_, err := testCase.g.DijkstraAll(testCase.src, testCase.dest)
		if err == nil {
			t.Errorf("expect error on %s", testCase.name)
		}
	}
}

var unweightedTestCases = []struct {
	g        *Graph
	src      VertexID
	dest     VertexID
	expected Path
}{
	{
		g:        NewGraph().LinkBoth(IntVertex(1), IntVertex(2), 1),
		src:      1,
		dest:     2,
		expected: Path{Stops: []Vertex{IntVertex(1), IntVertex(2)}, Weight: 1},
	},
	{
		g: NewGraph().
			LinkBoth(IntVertex(1), IntVertex(2), 1).
			LinkBoth(IntVertex(2), IntVertex(3), 1).
			LinkBoth(IntVertex(3), IntVertex(4), 1).
			LinkBoth(IntVertex(4), IntVertex(5), 1).
			LinkBoth(IntVertex(5), IntVertex(1), 1),
		src:      1,
		dest:     4,
		expected: Path{Stops: []Vertex{IntVertex(1), IntVertex(5), IntVertex(4)}, Weight: 2},
	},
}

func TestBFS(t *testing.T) {
	for _, testCase := range unweightedTestCases {
		actual, err := testCase.g.BFS(testCase.src, testCase.dest)
		if err != nil {
			t.Errorf("expected: %v, actual error: %s", testCase.expected, err)
		}
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("expected: %v, actual: %v", testCase.expected, actual)
		}
	}
}

var weightedTestCases = []struct {
	g        *Graph
	src      VertexID
	dest     VertexID
	expected []Path
}{
	{
		g:    NewGraph().LinkBoth(IntVertex(1), IntVertex(2), 1),
		src:  1,
		dest: 2,
		expected: []Path{
			Path{Stops: []Vertex{IntVertex(1), IntVertex(2)}, Weight: 1},
		},
	},
	{
		g: NewGraph().
			LinkBoth(IntVertex(1), IntVertex(2), 2).
			LinkBoth(IntVertex(2), IntVertex(3), 2).
			LinkBoth(IntVertex(3), IntVertex(4), 1).
			LinkBoth(IntVertex(4), IntVertex(5), 1).
			LinkBoth(IntVertex(5), IntVertex(1), 1),
		src:  1,
		dest: 3,
		expected: []Path{
			Path{Stops: []Vertex{IntVertex(1), IntVertex(5), IntVertex(4), IntVertex(3)}, Weight: 3},
			Path{Stops: []Vertex{IntVertex(1), IntVertex(2), IntVertex(3)}, Weight: 4},
		},
	},
	{
		g: NewGraph().
			LinkBoth(IntVertex(1), IntVertex(2), 5).
			LinkBoth(IntVertex(2), IntVertex(10), 5).
			LinkBoth(IntVertex(1), IntVertex(3), 1).
			LinkBoth(IntVertex(3), IntVertex(4), 1).
			LinkBoth(IntVertex(4), IntVertex(5), 1).
			LinkBoth(IntVertex(5), IntVertex(6), 1).
			LinkBoth(IntVertex(6), IntVertex(10), 1).
			LinkBoth(IntVertex(1), IntVertex(7), 2).
			LinkBoth(IntVertex(7), IntVertex(8), 2).
			LinkBoth(IntVertex(8), IntVertex(9), 2).
			LinkBoth(IntVertex(9), IntVertex(10), 2),
		src:  1,
		dest: 10,
		expected: []Path{
			Path{Stops: []Vertex{IntVertex(1), IntVertex(3), IntVertex(4), IntVertex(5), IntVertex(6), IntVertex(10)}, Weight: 5},
			Path{Stops: []Vertex{IntVertex(1), IntVertex(7), IntVertex(8), IntVertex(9), IntVertex(10)}, Weight: 8},
			Path{Stops: []Vertex{IntVertex(1), IntVertex(2), IntVertex(10)}, Weight: 10},
		},
	}}

func TestDijkstra(t *testing.T) {
	for _, testCase := range weightedTestCases {
		expected := testCase.expected[0]
		actual, err := testCase.g.Dijkstra(testCase.src, testCase.dest)
		if err != nil {
			t.Errorf("expected: %v, actual error: %s", expected, err)
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	}
}

func TestDijkstraAll(t *testing.T) {
	for _, testCase := range weightedTestCases {
		actual, err := testCase.g.DijkstraAll(testCase.src, testCase.dest)
		if err != nil {
			t.Errorf("expected: %v, actual error: %s", testCase.expected, err)
		}
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("expected: %v, actual: %v", testCase.expected, actual)
		}
	}
}

//// Benchmarks on path searching algorithms
func BenchmarkGraphBFS(b *testing.B) {
	var g = buildGraph(loadAllStations(), TravelCostByStop{})
	var source = StationID{line: "CC", number: 19}
	var destination = StationID{line: "DT", number: 15}

	for i := 0; i < b.N; i++ {
		g.BFS(source, destination)
	}
}

func BenchmarkGraphDijkstra(b *testing.B) {
	var g = buildGraph(loadAllStations(), TravelCostByStop{})
	var source = StationID{line: "CC", number: 19}
	var destination = StationID{line: "DT", number: 15}

	for i := 0; i < b.N; i++ {
		g.Dijkstra(source, destination)
	}
}

func BenchmarkGraphDijkstraAll(b *testing.B) {
	var g = buildGraph(loadAllStations(), TravelCostByStop{})
	var source = StationID{line: "CC", number: 19}
	var destination = StationID{line: "DT", number: 15}

	for i := 0; i < b.N; i++ {
		g.DijkstraAll(source, destination)
	}
}
