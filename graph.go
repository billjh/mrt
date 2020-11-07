package main

import (
	"errors"
	"math"
	"sort"
)

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
	Vertices map[VertexID]*Vertex
	Edges    map[VertexID]map[VertexID]Weight
}

// NewGraph creates an empty Graph, and returns its reference.
func NewGraph() *Graph {
	return &Graph{
		Vertices: make(map[VertexID]*Vertex),
		Edges:    make(map[VertexID]map[VertexID]Weight),
	}
}

// Add a stand-alone vertex to the Graph.
func (g *Graph) Add(v Vertex) *Graph {
	// add or replace reference to station
	g.Vertices[v.ID()] = &v
	// initialize edges map if not done so
	if g.Edges[v.ID()] == nil {
		g.Edges[v.ID()] = make(map[VertexID]Weight)
	}
	return g
}

// LinkBoth addes both vertices and the edges in bi-direction to the Graph
func (g *Graph) LinkBoth(v, u Vertex, w Weight) *Graph {
	g.Add(u)
	g.Add(v)
	g.Edges[u.ID()][v.ID()] = w
	g.Edges[v.ID()][u.ID()] = w
	return g
}

// Path records the stop and total weight from source to destination in a Graph
type Path struct {
	Stops  []Vertex
	Weight Weight
}

// ErrorSourceNotFound is returned by path-finding algorithms when source does not exist in graph.
var ErrorSourceNotFound = errors.New("source not found")

// ErrorDestinationNotFound is returned by path-finding algorithms when destination does not exist in graph.
var ErrorDestinationNotFound = errors.New("destination not found")

// ErrorSourceDestinationSame is returned by path-finding algorithms when source and destinatino is the same.
var ErrorSourceDestinationSame = errors.New("source and destination are the same")

// ErrorPathNotFound is returned by path-finding algorithms when no path exists.
var ErrorPathNotFound = errors.New("path not found")

// UnweightedSearch use Graph.BFS when all is false; otherwise use
// Graph.DijkstraAll
func (g *Graph) UnweightedSearch(src, dest VertexID, all bool) ([]Path, error) {
	if all {
		return g.DijkstraAll(src, dest)
	}
	p, err := g.BFS(src, dest)
	return []Path{p}, err
}

// WeightedSearch use Graph.Dijkstra when all is false; otherwise use
// Graph.DijkstraAll
func (g *Graph) WeightedSearch(src, dest VertexID, all bool) ([]Path, error) {
	if all {
		return g.DijkstraAll(src, dest)
	}
	p, err := g.Dijkstra(src, dest)
	return []Path{p}, err
}

// BFS finds the shortest path from source to destination and ignores edge weights.
// It returns error when
// 1) source or destination does not exist in the Graph;
// 2) source and destination are the same;
// 3) no path is found.
func (g *Graph) BFS(src, dest VertexID) (Path, error) {
	if err := validate(g, src, dest); err != nil {
		return Path{}, err
	}
	parent := make(map[VertexID]VertexID)
	visited := map[VertexID]bool{src: true}
	queue := []VertexID{src}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current == dest {
			path := g.backtrack(current, parent)
			return Path{Stops: path, Weight: Weight(len(path) - 1)}, nil
		}
		for neighbor := range g.Edges[current] {
			if !visited[neighbor] {
				parent[neighbor] = current
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
	return Path{}, ErrorPathNotFound
}

// Dijkstra finds the path with minimum weight from source to destination in a Graph.
// It returns error when
// 1) source or destination does not exist in the Graph;
// 2) source and destination are the same;
// 3) no path is found.
func (g *Graph) Dijkstra(src, dest VertexID) (Path, error) {
	if err := validate(g, src, dest); err != nil {
		return Path{}, err
	}
	parent := make(map[VertexID]VertexID)
	visited := map[VertexID]bool{src: true}
	dist := map[VertexID]Weight{src: 0}

	for len(dist) > 0 {
		// pop the nearest vertex
		current := minDist(dist)
		currentWeight := dist[current]
		delete(dist, current)
		visited[current] = true

		if current == dest {
			p := g.backtrack(current, parent)
			return Path{
				Stops:  p,
				Weight: currentWeight,
			}, nil
		}

		for neighbor, edgeWeight := range g.Edges[current] {
			if !visited[neighbor] {
				alt := currentWeight + edgeWeight
				neighborWeight, ok := dist[neighbor]
				if !ok || alt < neighborWeight {
					dist[neighbor] = alt
					parent[neighbor] = current
				}
			}
		}
	}
	return Path{}, ErrorPathNotFound
}

// DijkstraAll finds all the paths from source to destination in a Graph and sorts them by
// total weight in descending order. It returns error when
// 1) source or destination does not exist in the Graph;
// 2) source and destination are the same;
// 3) no path is found.
func (g *Graph) DijkstraAll(src, dest VertexID) ([]Path, error) {
	if err := validate(g, src, dest); err != nil {
		return nil, err
	}

	paths := []Path{}

	parent := make(map[VertexID]VertexID)
	visited := map[VertexID]bool{src: true}
	dist := map[VertexID]Weight{src: 0}

	for len(dist) > 0 {
		// pop the nearest vertex
		current := minDist(dist)
		currentWeight := dist[current]
		delete(dist, current)
		visited[current] = true

		if current == dest {
			// only stop when all neighbors of dest have been visited
			hasUnvisited := false
			for neighbor := range g.Edges[current] {
				if !visited[neighbor] {
					hasUnvisited = true
					break
				}
			}
			if !hasUnvisited {
				break
			}
		}

		for neighbor, edgeWeight := range g.Edges[current] {
			// do not explore neighbors of destination
			if current == dest {
				break
			}

			// record down the all the paths
			if neighbor == dest {
				p := append(g.backtrack(current, parent), *g.Vertices[dest])
				paths = append(paths, Path{
					Stops:  p,
					Weight: currentWeight + edgeWeight,
				})
			}

			if !visited[neighbor] {
				alt := currentWeight + edgeWeight
				neighborWeight, ok := dist[neighbor]
				if !ok || alt < neighborWeight {
					dist[neighbor] = alt
					parent[neighbor] = current
				}
			}
		}
	}
	// returns error if no path is found
	if len(paths) == 0 {
		return nil, ErrorPathNotFound
	}
	// sort the paths by weight in descending order
	sort.Slice(paths, func(i, j int) bool { return paths[i].Weight < paths[j].Weight })
	return paths, nil
}

// validate the source and destination for path finding algorithms
func validate(g *Graph, src, dest VertexID) error {
	if _, ok := g.Vertices[src]; !ok {
		return ErrorSourceNotFound
	}
	if _, ok := g.Vertices[dest]; !ok {
		return ErrorDestinationNotFound
	}
	if src == dest {
		return ErrorSourceDestinationSame
	}
	return nil
}

// backtrack is a helper function which constructs the path with parent map
func (g *Graph) backtrack(current VertexID, parent map[VertexID]VertexID) []Vertex {
	path := []Vertex{*g.Vertices[current]}
	for {
		if p, ok := parent[current]; ok {
			path = append([]Vertex{*g.Vertices[p]}, path...)
			current = p
		} else {
			return path
		}
	}
}

// minDist is a helper function which finds the nearest unvisited vertex,
// and it assumes non empty dist map
func minDist(dist map[VertexID]Weight) VertexID {
	var min VertexID
	minWeight := math.MaxInt32
	for v, w := range dist {
		if int(w) < minWeight {
			min = v
			minWeight = int(w)
		}
	}
	return min
}
