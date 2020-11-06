package main

import (
	"reflect"
	"testing"
	"time"

	"github.com/billjh/zendesk-mrt/graph"
)

// Peak hours (6am-9am and 6pm-9pm on Mon-Fri)
// Night hours (10pm-6am on Mon-Sun)
// Non-Peak hours (all other times)
func TestTravelPeriod(t *testing.T) {
	peak := "peak"
	night := "night"
	nonpeak := "nonpeak"

	timeLayout := "2006-01-02T15:04"
	sun := "2020-11-08T"
	mon := "2020-11-09T"
	fri := "2020-11-13T"
	sat := "2020-11-14T"
	for _, testCase := range []struct {
		t string
		h string
	}{
		{t: mon + "05:59", h: night},
		{t: fri + "05:59", h: night},
		{t: sat + "05:59", h: night},
		{t: sun + "05:59", h: night},

		{t: mon + "06:01", h: peak},
		{t: fri + "06:01", h: peak},
		{t: sat + "06:01", h: nonpeak},
		{t: sun + "06:01", h: nonpeak},

		{t: mon + "08:59", h: peak},
		{t: fri + "08:59", h: peak},
		{t: sat + "08:59", h: nonpeak},
		{t: sun + "08:59", h: nonpeak},

		{t: mon + "09:01", h: nonpeak},
		{t: fri + "09:01", h: nonpeak},
		{t: sat + "09:01", h: nonpeak},
		{t: sun + "09:01", h: nonpeak},

		{t: mon + "17:59", h: nonpeak},
		{t: fri + "17:59", h: nonpeak},
		{t: sat + "17:59", h: nonpeak},
		{t: sun + "17:59", h: nonpeak},

		{t: mon + "18:01", h: peak},
		{t: fri + "18:01", h: peak},
		{t: sat + "18:01", h: nonpeak},
		{t: sun + "18:01", h: nonpeak},

		{t: mon + "20:59", h: peak},
		{t: fri + "20:59", h: peak},
		{t: sat + "20:59", h: nonpeak},
		{t: sun + "20:59", h: nonpeak},

		{t: mon + "21:01", h: nonpeak},
		{t: fri + "21:01", h: nonpeak},
		{t: sat + "21:01", h: nonpeak},
		{t: sun + "21:01", h: nonpeak},

		{t: mon + "21:59", h: nonpeak},
		{t: fri + "21:59", h: nonpeak},
		{t: sat + "21:59", h: nonpeak},
		{t: sun + "21:59", h: nonpeak},

		{t: mon + "22:01", h: night},
		{t: fri + "22:01", h: night},
		{t: sat + "22:01", h: night},
		{t: sun + "22:01", h: night},
	} {
		p, err := time.Parse(timeLayout, testCase.t)
		if err != nil {
			t.Error(err)
		}
		var actual string
		switch {
		case isPeakHours(p):
			actual = peak
		case isNightHours(p):
			actual = night
		default:
			actual = nonpeak
		}
		if actual != testCase.h {
			t.Errorf("%s %s expected: %s, actual: %s", p.Weekday(), testCase.t, testCase.h, actual)
		}
	}
}

func TestBuildGraph(t *testing.T) {
	travelCost := TravelCostByTime{
		interchange: 1,
		lines: map[string]graph.Weight{
			"NE": 3,
		},
		lineDefault: 2,
	}
	stations := []Station{
		// NE5,Clarke Quay
		// NE6,Dhoby Ghaut
		// CC1,Dhoby Ghaut
		// CC2,Bras Basah
		Station{
			id:   StationID{line: "NE", number: 5},
			name: "Clarke Quay",
		},
		Station{
			id:   StationID{line: "NE", number: 6},
			name: "Dhoby Ghaut",
		},
		Station{
			id:   StationID{line: "CC", number: 1},
			name: "Dhoby Ghaut",
		},
		Station{
			id:   StationID{line: "CC", number: 2},
			name: "Bras Basah",
		},
	}
	g := buildGraph(stations, travelCost)

	// need to use the type same as graph.Graph.Edges
	expectedEdges := map[graph.VertexID]map[graph.VertexID]graph.Weight{
		StationID{line: "NE", number: 5}: map[graph.VertexID]graph.Weight{
			StationID{line: "NE", number: 6}: 3,
		},
		StationID{line: "NE", number: 6}: map[graph.VertexID]graph.Weight{
			StationID{line: "NE", number: 5}: 3,
			StationID{line: "CC", number: 1}: 1,
		},
		StationID{line: "CC", number: 1}: map[graph.VertexID]graph.Weight{
			StationID{line: "NE", number: 6}: 1,
			StationID{line: "CC", number: 2}: 2,
		},
		StationID{line: "CC", number: 2}: map[graph.VertexID]graph.Weight{
			StationID{line: "CC", number: 1}: 2,
		},
	}
	if !reflect.DeepEqual(g.Edges, expectedEdges) {
		t.Errorf("Edges not match\nexpected: %v\n  actual: %v", expectedEdges, g.Edges)
	}
}

func TestByStops(t *testing.T) {
	for _, testCase := range []struct {
		src      StationID
		dest     StationID
		expected []string
	}{
		{
			src:      StationID{line: "CC", number: 21},
			dest:     StationID{line: "DT", number: 14},
			expected: []string{"CC21", "CC20", "CC19", "DT9", "DT10", "DT11", "DT12", "DT13", "DT14"},
		},
	} {
		path, err := NewNavigator().ByStops(testCase.src, testCase.dest)
		if err != nil {
			t.Error(err)
		}
		actual := []string{}
		for _, s := range path {
			actual = append(actual, s.(Station).id.String())
		}
		if !reflect.DeepEqual(testCase.expected, actual) {
			t.Errorf("expected: %v, actual: %v", testCase.expected, actual)
		}
	}
}

func TestByTime(t *testing.T) {
	for _, testCase := range []struct {
		src      StationID
		dest     StationID
		timeStr  string
		expected []string
		weight   int
	}{
		{
			src:      StationID{line: "EW", number: 27},
			dest:     StationID{line: "DT", number: 12},
			timeStr:  "2020-11-09T06:01",
			expected: []string{"EW27", "EW26", "EW25", "EW24", "EW23", "EW22", "EW21", "CC22", "CC21", "CC20", "CC19", "DT9", "DT10", "DT11", "DT12"},
			weight:   150,
		},
	} {
		travelTime, err := time.Parse("2006-01-02T15:04", testCase.timeStr)
		if err != nil {
			t.Error(err)
		}
		path, err := NewNavigator().ByTime(testCase.src, testCase.dest, travelTime)
		if err != nil {
			t.Error(err)
		}
		actual := []string{}
		for _, s := range path.Path {
			actual = append(actual, s.(Station).id.String())
		}
		if !reflect.DeepEqual(testCase.expected, actual) {
			t.Errorf("expected: %v, actual: %v", testCase.expected, actual)
		}
	}
}
