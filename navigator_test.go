package main

import (
	"reflect"
	"testing"
	"time"
)

func TestBuildGraph(t *testing.T) {
	travelCost := TravelCostByTime{
		interchange: 1,
		lines: map[string]Weight{
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

	// need to use the type same as Graph.Edges
	expectedEdges := map[VertexID]map[VertexID]Weight{
		StationID{line: "NE", number: 5}: map[VertexID]Weight{
			StationID{line: "NE", number: 6}: 3,
		},
		StationID{line: "NE", number: 6}: map[VertexID]Weight{
			StationID{line: "NE", number: 5}: 3,
			StationID{line: "CC", number: 1}: 1,
		},
		StationID{line: "CC", number: 1}: map[VertexID]Weight{
			StationID{line: "NE", number: 6}: 1,
			StationID{line: "CC", number: 2}: 2,
		},
		StationID{line: "CC", number: 2}: map[VertexID]Weight{
			StationID{line: "CC", number: 1}: 2,
		},
	}
	if !reflect.DeepEqual(g.Edges, expectedEdges) {
		t.Errorf("Edges not match\nexpected: %v\n  actual: %v", expectedEdges, g.Edges)
	}
}

type ExpectedPath struct {
	path   []string
	weight Weight
}

func TestNavigateByStops(t *testing.T) {
	for _, testCase := range []struct {
		src         string
		dest        string
		all         bool
		expectError bool
		expected    []ExpectedPath
	}{
		{
			src:         "???",
			dest:        "Lakeside",
			expectError: true,
		},
		{
			src:         "Lakeside",
			dest:        "???",
			expectError: true,
		},
		{
			src:         "Chinatown",
			dest:        "Chinatown",
			expectError: true,
		},
		{
			src:         "NE4",
			dest:        "Chinatown",
			expectError: true,
		},
		{
			src:  "NE4",
			dest: "DT19",
			all:  false,
			expected: []ExpectedPath{
				ExpectedPath{weight: 1, path: []string{"NE4", "DT19"}},
			},
		},
		{
			src:  "CC21",
			dest: "DT14",
			all:  false,
			expected: []ExpectedPath{
				ExpectedPath{weight: 8, path: []string{"CC21", "CC20", "CC19", "DT9", "DT10", "DT11", "DT12", "DT13", "DT14"}},
			},
		},
		{
			src:  "Jurong East",
			dest: "HarbourFront",
			all:  true,
			expected: []ExpectedPath{
				ExpectedPath{weight: 10, path: []string{"EW24", "EW23", "EW22", "EW21", "EW20", "EW19", "EW18", "EW17", "EW16", "NE3", "NE1"}},
				ExpectedPath{weight: 11, path: []string{"EW24", "EW23", "EW22", "EW21", "CC22", "CC23", "CC24", "CC25", "CC26", "CC27", "CC28", "CC29"}},
			},
		},
	} {
		paths, err := NewNavigator().NavigateByStops(testCase.src, testCase.dest, testCase.all)
		if testCase.expectError {
			if err == nil {
				t.Errorf("expect error '%s' to '%s'", testCase.src, testCase.dest)
			}
		} else {
			actual := []ExpectedPath{}
			for _, p := range paths {
				actual = append(actual, ExpectedPath{path: pathToStringSlice(p.Stops), weight: p.Weight})
			}
			if !reflect.DeepEqual(testCase.expected, actual) {
				t.Errorf("\nexpected: %v, \n  actual: %v", testCase.expected, actual)
			}
		}
	}
}

func TestNavigateByTime(t *testing.T) {
	peakHours := "2020-11-09T06:01"
	nightHours := "2020-11-09T05:59"
	nonPeakHours := "2020-11-08T06:01"

	for _, testCase := range []struct {
		src         string
		dest        string
		timeStr     string
		all         bool
		expectError bool
		expected    []ExpectedPath
	}{
		{
			src:         "???",
			dest:        "Lakeside",
			timeStr:     peakHours,
			expectError: true,
		},
		{
			src:         "Lakeside",
			dest:        "???",
			timeStr:     peakHours,
			expectError: true,
		},
		{
			src:         "Chinatown",
			dest:        "Chinatown",
			timeStr:     peakHours,
			expectError: true,
		},
		{
			src:         "NE4",
			dest:        "Chinatown",
			timeStr:     peakHours,
			expectError: true,
		},
		{
			src:     "EW27",
			dest:    "DT12",
			timeStr: peakHours,
			all:     false,
			expected: []ExpectedPath{
				ExpectedPath{weight: 150, path: []string{"EW27", "EW26", "EW25", "EW24", "EW23", "EW22", "EW21", "CC22", "CC21", "CC20", "CC19", "DT9", "DT10", "DT11", "DT12"}},
			},
		},
		{
			src:     "CC19",
			dest:    "CC4",
			timeStr: nonPeakHours,
			all:     false,
			expected: []ExpectedPath{
				ExpectedPath{weight: 68, path: []string{"CC19", "DT9", "DT10", "DT11", "DT12", "DT13", "DT14", "DT15", "CC4"}},
			},
		},
		{
			src:     "CC19",
			dest:    "CC4",
			timeStr: nightHours,
			all:     false,
			expected: []ExpectedPath{
				ExpectedPath{weight: 140, path: []string{"CC19", "CC17", "CC16", "CC15", "CC14", "CC13", "CC12", "CC11", "CC10", "CC9", "CC8", "CC7", "CC6", "CC5", "CC4"}},
			},
		},
		{
			src:     "Boon Lay",
			dest:    "Little India",
			timeStr: peakHours,
			all:     false,
			expected: []ExpectedPath{
				ExpectedPath{weight: 150, path: []string{"EW27", "EW26", "EW25", "EW24", "EW23", "EW22", "EW21", "CC22", "CC21", "CC20", "CC19", "DT9", "DT10", "DT11", "DT12"}},
			},
		},
		{
			src:     "Jurong East",
			dest:    "HarbourFront",
			timeStr: peakHours,
			all:     true,
			expected: []ExpectedPath{
				ExpectedPath{weight: 107, path: []string{"EW24", "EW23", "EW22", "EW21", "EW20", "EW19", "EW18", "EW17", "EW16", "NE3", "NE1"}},
				ExpectedPath{weight: 115, path: []string{"EW24", "EW23", "EW22", "EW21", "CC22", "CC23", "CC24", "CC25", "CC26", "CC27", "CC28", "CC29"}},
			},
		},
	} {
		travelTime, err := time.Parse("2006-01-02T15:04", testCase.timeStr)
		if err != nil {
			t.Error(err)
		}
		paths, err := NewNavigator().NavigateByTime(testCase.src, testCase.dest, travelTime, testCase.all)
		if testCase.expectError {
			if err == nil {
				t.Errorf("expect error '%s' to '%s'", testCase.src, testCase.dest)
			}
		} else {
			actual := []ExpectedPath{}
			for _, p := range paths {
				actual = append(actual, ExpectedPath{path: pathToStringSlice(p.Stops), weight: p.Weight})
			}
			if !reflect.DeepEqual(testCase.expected, actual) {
				t.Errorf("\nexpected: %v, \n  actual: %v", testCase.expected, actual)
			}
		}
	}
}

// pathToStringSlice is a helper function convert Path to station codes in string
func pathToStringSlice(path []Vertex) []string {
	actual := []string{}
	for _, s := range path {
		actual = append(actual, s.(Station).id.String())
	}
	return actual
}

//// Benchmarks on Navigator methods
func BenchmarkNavigateByStopsSingle(b *testing.B) {
	var navigatorForBenchmark = NewNavigator()
	var source, destination = "Botanic Garden", "Promenade"

	for i := 0; i < b.N; i++ {
		navigatorForBenchmark.NavigateByStops(source, destination, false)
	}
}

func BenchmarkNavigateByStopsAll(b *testing.B) {
	var navigatorForBenchmark = NewNavigator()
	var source, destination = "Botanic Garden", "Promenade"

	for i := 0; i < b.N; i++ {
		navigatorForBenchmark.NavigateByStops(source, destination, true)
	}
}

func BenchmarkNavigateByTimeSingle(b *testing.B) {
	var navigatorForBenchmark = NewNavigator()
	var source, destination = "Botanic Garden", "Promenade"
	var travelTime, _ = time.Parse("2006-01-02T15:04", "2020-11-09T06:01")

	for i := 0; i < b.N; i++ {
		navigatorForBenchmark.NavigateByTime(source, destination, travelTime, false)
	}
}

func BenchmarkNavigateByTimeAll(b *testing.B) {
	var navigatorForBenchmark = NewNavigator()
	var source, destination = "Botanic Garden", "Promenade"
	var travelTime, _ = time.Parse("2006-01-02T15:04", "2020-11-09T06:01")

	for i := 0; i < b.N; i++ {
		navigatorForBenchmark.NavigateByTime(source, destination, travelTime, true)
	}
}
