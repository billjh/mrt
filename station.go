package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// StationCode is the code for a MRT station in Singapore
// It consists of two fields, 2-alphabet MRT line and int station number
// For example, StationCode{line: "EW", number: 14}
type StationCode struct {
	line   string
	number int
}

// Station represents a MRT station in Singapore
type Station struct {
	code        StationCode
	name        string
	openingDate time.Time
}

// NewStationCode constructs a StationCode from string
// and returns error on invalid format
func NewStationCode(code string) (StationCode, error) {
	matched, err := regexp.MatchString(`^[a-zA-Z]{2}\d{1,2}$`, code)
	if err != nil {
		return StationCode{}, err
	}
	if !matched {
		return StationCode{}, fmt.Errorf("invalid station code %s", code)
	}
	number, err := strconv.Atoi(code[2:])
	if err != nil {
		// this shouldn't happen though
		return StationCode{}, err
	}
	return StationCode{line: strings.ToUpper(code[:2]), number: number}, nil
}
