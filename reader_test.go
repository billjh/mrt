package main

import (
	"strings"
	"testing"
	"time"
)

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
					code:        "EW23",
					name:        "Clementi",
					openingDate: time.Date(1988, 3, 12, 0, 0, 0, 0, time.UTC),
				},
				Station{
					code:        "EW24",
					name:        "Jurong East",
					openingDate: time.Date(1988, 11, 5, 0, 0, 0, 0, time.UTC),
				},
				Station{
					code:        "EW25",
					name:        "Chinese Garden",
					openingDate: time.Date(1988, 11, 5, 0, 0, 0, 0, time.UTC),
				},
				Station{
					code:        "EW26",
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
