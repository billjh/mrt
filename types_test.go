package main

import (
	"testing"
)

func TestMRTLine(t *testing.T) {
	for _, testCase := range []struct {
		station  Station
		expected string
	}{
		{
			station:  Station{},
			expected: "",
		},
		{
			station:  Station{code: "TE1"},
			expected: "TE",
		},
		{
			station:  Station{code: "DT19"},
			expected: "DT",
		},
	} {
		actual := testCase.station.MRTLine()
		if actual != testCase.expected {
			t.Errorf("expect: %s, actual: %s", testCase.expected, actual)
		}
	}
}
