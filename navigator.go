package main

import (
	"sort"
	"time"
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

// NavigateByStops returns shortest paths between two Stations or any error encountered.
// It accepts source and destination input as string, which can be either StationID like "DT1"
// or station name like "Bukit Panjang". If all is set to true, all the paths ordered by
// number of stops are returned instead just the shortest.
func (n *Navigator) NavigateByStops(srcStr, destStr string, all bool) ([]Path, error) {
	allSrc, srcIsID, err := searchStations(n.allStations, srcStr)
	if err != nil {
		return nil, ErrorSourceNotFound
	}
	allDest, destIsID, err := searchStations(n.allStations, destStr)
	if err != nil {
		return nil, ErrorDestinationNotFound
	}

	g := buildGraph(n.allStations, TravelCostByStop{})

	paths := []Path{}

	for _, src := range allSrc {
		for _, dest := range allDest {
			ps, err := g.UnweightedSearch(src, dest, all)
			if err != nil {
				continue
			}
			for _, p := range ps {
				l := len(p.Stops)
				if l < 2 {
					continue
				}
				// filter out paths that start interchanging when source is not pinned to an ID
				if !srcIsID && p.Stops[0].(Station).name == p.Stops[1].(Station).name {
					continue
				}
				// filter out paths that end interchanging when destination is not pinned to an ID
				if !destIsID && p.Stops[l-1].(Station).name == p.Stops[l-2].(Station).name {
					continue
				}
				paths = append(paths, p)
			}
		}
	}

	if len(paths) == 0 {
		return nil, ErrorPathNotFound
	}

	sort.Slice(paths, func(i, j int) bool { return paths[i].Weight < paths[j].Weight })

	if all {
		return paths, nil
	}
	return paths[:1], nil
}

// NavigateByTime returns fastest paths between two Stations or any error encountered, knowing the
// time of travel.
// It accepts source and destination input as string, which can be either StationID like "DT1"
// or station name like "Bukit Panjang". If all is set to true, all paths ordered by
// estimated time are returned instead just the fastest.
func (n *Navigator) NavigateByTime(srcStr, destStr string, t time.Time, all bool) ([]Path, error) {
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

	allSrc, srcIsID, err := searchStations(openingStations, srcStr)
	if err != nil {
		return nil, ErrorSourceNotFound
	}
	allDest, destIsID, err := searchStations(openingStations, destStr)
	if err != nil {
		return nil, ErrorDestinationNotFound
	}

	g := buildGraph(openingStations, getTravelCostByTime(t))

	paths := []Path{}

	for _, src := range allSrc {
		for _, dest := range allDest {
			ps, err := g.WeightedSearch(src, dest, all)
			if err != nil {
				continue
			}
			for _, p := range ps {
				l := len(p.Stops)
				if l < 2 {
					continue
				}
				// filter out paths that start interchanging when source is not pinned to an ID
				if !srcIsID && p.Stops[0].(Station).name == p.Stops[1].(Station).name {
					continue
				}
				// filter out paths that end interchanging when destination is not pinned to an ID
				if !destIsID && p.Stops[l-1].(Station).name == p.Stops[l-2].(Station).name {
					continue
				}
				paths = append(paths, p)
			}
		}
	}

	if len(paths) == 0 {
		return nil, ErrorPathNotFound
	}

	sort.Slice(paths, func(i, j int) bool { return paths[i].Weight < paths[j].Weight })

	if all {
		return paths, nil
	}
	return paths[:1], nil
}

// buildGraph takes a list of Stations and connects them in a Graph:
//
// 1) each Station on the same MRT line is connected to its adjacent Stations,
// eg. EW1 <-> EW2 <-> EW4 (assuming EW3 not exists)
//
// 2) Stations with the same name but different StationIDs will be treated as
// interchange stations, so they will be connected to each other
// eg. NS24 Dhoby Ghaut <-> CC1 Dhoby Ghaut <-> NE6 Dhoby Ghaut (<-> NS24 Dhoby Ghaut)
func buildGraph(stations []Station, cost TravelCost) *Graph {
	g := NewGraph()

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
