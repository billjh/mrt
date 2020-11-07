package main

import (
	"sort"
	"time"

	"github.com/billjh/zendesk-mrt/graph"
)

// Navigator holds a map of all Stations and provides multiple navigating methods
type Navigator struct {
	allStations []Station
}

// NewNavigator loads all Stations and returns a Navigator instance
func NewNavigator() *Navigator {
	return &Navigator{
		allStations: loadAllStations(),
	}
}

// ByStops gives the shortest path by number of stops, or any error encountered
func (n *Navigator) ByStops(src, dest StationID) (graph.Path, error) {
	g := buildGraph(n.allStations, TravelCostByStop{})
	return g.BFS(src, dest)
}

// ByTime gives the fatest pathes by time taken, or any error encountered, knowing the time of travel.
// If all is set false, return only the fastest path; otherwise return all pathes ordered by time taken.
func (n *Navigator) ByTime(src, dest StationID, t time.Time, all bool) ([]graph.Path, error) {
	// get opening stations at the time of travel
	openingStations := []Station{}
	for _, station := range n.allStations {
		// remove if travel before the station exists
		if t.Before(station.openingDate) {
			continue
		}
		// DT, CG and CE lines do not operate at night
		if isNightHours(t) && stopAtNight(station.id.line) {
			continue
		}
		openingStations = append(openingStations, station)
	}

	g := buildGraph(openingStations, getTravelCostByTime(t))

	// use Graph.DijkstraAll to find all pathes if all flag is set
	if all {
		return g.DijkstraAll(src, dest)
	}

	// otherwise, use Graph.Dijkstra to get only the fatest
	p, err := g.Dijkstra(src, dest)

	return []graph.Path{p}, err
}

// buildGraph takes a list of Stations and connects them in a Graph:
//
// 1) each Station on the same MRT line is connected to its adjacent Stations,
// eg. EW1 <-> EW2 <-> EW4 (assuming EW3 not exists)
//
// 2) Stations with the same name but different StationIDs will be treated as
// interchange stations, so they will be connected to each other
// eg. NS24 Dhoby Ghaut <-> CC1 Dhoby Ghaut <-> NE6 Dhoby Ghaut (<-> NS24 Dhoby Ghaut)
func buildGraph(stations []Station, cost TravelCost) *graph.Graph {
	g := graph.NewGraph()

	// add all Stations as graph vertices
	for _, s := range stations {
		g.Add(s)
	}

	// link adajent Stations on the same line
	for line, ss := range groupBy(stations, func(s Station) string { return s.id.line }) {
		sort.Slice(ss, func(i, j int) bool { return ss[i].id.number < ss[j].id.number })
		for i := 1; i < len(ss); i++ {
			g.LinkBoth(ss[i-1], ss[i], cost.OnLine(line))
		}
	}

	// link interchange Stations
	for _, ss := range groupBy(stations, func(s Station) string { return s.name }) {
		for i := 0; i < len(ss); i++ {
			for j := i + 1; j < len(ss); j++ {
				g.LinkBoth(ss[i], ss[j], cost.Interchange())
			}
		}
	}

	return g
}

// groupBy is a helper function to group Stations by key func
func groupBy(stations []Station, key func(Station) string) map[string][]Station {
	m := make(map[string][]Station)
	for _, s := range stations {
		if m[key(s)] == nil {
			m[key(s)] = []Station{s}
		} else {
			m[key(s)] = append(m[key(s)], s)
		}
	}
	return m
}
