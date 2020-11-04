package main

// VertexID is a generic interface type to represent vertex's identity.
type VertexID interface{}

// Vertex must implement this interface to be used in Graph.
type Vertex interface {
	ID() VertexID
}

// Weight is assumed to be int. It's used for edge weight and path weight.
type Weight int

// Graph contains all the Vertex references by a map accessed by ID.
// It also stores the edges for each Vertex, which is a map of weights.
type Graph struct {
	vertices map[VertexID]*Vertex
	edges    map[VertexID]map[VertexID]Weight
}

// NewGraph creates an empty graph, and returns its reference.
func NewGraph() *Graph {
	return &Graph{
		vertices: make(map[VertexID]*Vertex),
		edges:    make(map[VertexID]map[VertexID]Weight),
	}
}

// Add a stand-alone vertex to the graph.
func (g *Graph) Add(v Vertex) {
	// add or replace reference to station
	g.vertices[v.ID()] = &v
	// initialize edges map if not done so
	if g.edges[v.ID()] == nil {
		g.edges[v.ID()] = make(map[VertexID]Weight)
	}
}

// LinkBoth addes both vertices and the edges in bi-direction to the graph
func (g *Graph) LinkBoth(v, u Vertex, w Weight) {
	g.Add(u)
	g.Add(v)
	g.edges[u.ID()][v.ID()] = w
	g.edges[v.ID()][u.ID()] = w
}

// Path records the stops and total weight in a graph from source to destination
type Path struct {
	stops  []VertexID
	weight Weight
}
