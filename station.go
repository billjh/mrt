package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// StationCode is the code for a MRT station in Singapore.
// It consists of two fields, 2-alphabet MRT line and int station number.
// For example, StationCode{line: "EW", number: 14}
type StationCode struct {
	line   string
	number int
}

// Station represents a MRT station in Singapore.
type Station struct {
	code        StationCode
	name        string
	openingDate time.Time
}

// NewStationCode constructs a StationCode from string
// and returns error on invalid format.
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

// ReadStations reads the stations from the given io.Reader.
// It assumes the format being:
//
// Station Code,Station Name,Opening Date
// EW23,Clementi,12 March 1988
// EW24,Jurong East,5 November 1988
// EW25,Chinese Garden,5 November 1988
// EW26,Lakeside,5 November 1988
//
func ReadStations(r io.Reader) ([]Station, error) {
	csvReader := csv.NewReader(r)

	// skip header row
	_, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	// parse each row and store in the final result array
	final := []Station{}

	// the format used by time.Parse function
	const openingDateFormat string = "2 January 2006"

	for _, record := range records {
		if len(record) != 3 {
			return nil, fmt.Errorf("record lenth not 3: %v", record)
		}
		code, err := NewStationCode(record[0])
		if err != nil {
			return nil, err
		}
		openingDate, err := time.Parse(openingDateFormat, record[2])
		if err != nil {
			return nil, err
		}
		final = append(final, Station{
			code:        code,
			name:        record[1],
			openingDate: openingDate,
		})
	}

	return final, nil
}
