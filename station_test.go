package main

import (
	"strings"
	"testing"
	"time"
)

func TestNewStationIDSuccess(t *testing.T) {
	for _, testCase := range []struct {
		input    string
		expected StationID
	}{
		{
			input:    "TE1",
			expected: StationID{line: "TE", number: 1},
		},
		{
			input:    "te1",
			expected: StationID{line: "TE", number: 1},
		},
		{
			input:    "DT19",
			expected: StationID{line: "DT", number: 19},
		},
	} {
		actual, err := NewStationID(testCase.input)
		if err != nil {
			t.Error(err)
		}
		if testCase.expected != actual {
			t.Errorf("expected: %v, actual: %v", testCase.expected, actual)
		}
	}
}

func TestNewStationIDError(t *testing.T) {
	for _, testCase := range []string{
		"",
		"DT",
		"1",
		"19",
		"1DT",
		"ABC1",
		"DT888",
	} {
		_, err := NewStationID(testCase)
		if err == nil {
			t.Errorf("expect error for input: %s", testCase)
		}
	}
}

func TestReadStations(t *testing.T) {
	for _, testCase := range []struct {
		fileContent string
		expected    []Station
	}{
		{
			fileContent: "Station Code,Station Name,Opening Date\n" +
				"EW23,Clementi,12 March 1988\n" +
				"EW24,Jurong East,5 November 1988\n" +
				"EW25,Chinese Garden,5 November 1988\n" +
				"EW26,Lakeside,5 November 1988\n",
			expected: []Station{
				Station{
					id:          StationID{line: "EW", number: 23},
					name:        "Clementi",
					openingDate: time.Date(1988, 3, 12, 0, 0, 0, 0, time.UTC),
				},
				Station{
					id:          StationID{line: "EW", number: 24},
					name:        "Jurong East",
					openingDate: time.Date(1988, 11, 5, 0, 0, 0, 0, time.UTC),
				},
				Station{
					id:          StationID{line: "EW", number: 25},
					name:        "Chinese Garden",
					openingDate: time.Date(1988, 11, 5, 0, 0, 0, 0, time.UTC),
				},
				Station{
					id:          StationID{line: "EW", number: 26},
					name:        "Lakeside",
					openingDate: time.Date(1988, 11, 5, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	} {
		s, err := ReadStations(strings.NewReader(testCase.fileContent))
		if err != nil {
			t.Error(err)
		}
		if len(s) != len(testCase.expected) {
			t.Errorf("len not match, expected: %d, actual: %d", len(testCase.expected), len(s))
		}
		for i, actual := range s {
			if actual != testCase.expected[i] {
				t.Errorf("item not match, expected: %v, actual: %v", testCase.expected[i], actual)
			}
		}
	}
}
