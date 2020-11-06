package main

import (
	"testing"
)

func TestNewStationCodeSuccess(t *testing.T) {
	for _, testCase := range []struct {
		input    string
		expected StationCode
	}{
		{
			input:    "TE1",
			expected: StationCode{line: "TE", number: 1},
		},
		{
			input:    "te1",
			expected: StationCode{line: "TE", number: 1},
		},
		{
			input:    "DT19",
			expected: StationCode{line: "DT", number: 19},
		},
	} {
		actual, err := NewStationCode(testCase.input)
		if err != nil {
			t.Error(err)
		}
		if testCase.expected != actual {
			t.Errorf("expected: %v, actual: %v", testCase.expected, actual)
		}
	}
}

func TestNewStationCodeError(t *testing.T) {
	for _, testCase := range []string{
		"",
		"DT",
		"1",
		"19",
		"1DT",
		"ABC1",
		"DT888",
	} {
		_, err := NewStationCode(testCase)
		if err == nil {
			t.Errorf("expect error for input: %s", testCase)
		}
	}
}
