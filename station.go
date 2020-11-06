package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/billjh/zendesk-mrt/graph"
)

// StationID is the code for a MRT station in Singapore.
// It consists of 2-letter line code and number.
// For example, StationID{line: "EW", number: 14}
type StationID struct {
	line   string
	number int
}

// Station represents a MRT station in Singapore.
type Station struct {
	id          StationID
	name        string
	openingDate time.Time
}

// ID implements graph.Vertex interface
func (s Station) ID() graph.VertexID {
	return s.id
}

// NewStationID constructs a StationID from string
// and returns error on invalid format.
func NewStationID(id string) (StationID, error) {
	matched, err := regexp.MatchString(`^[a-zA-Z]{2}\d{1,2}$`, id)
	if err != nil {
		return StationID{}, err
	}
	if !matched {
		return StationID{}, fmt.Errorf("invalid station id %s", id)
	}
	number, err := strconv.Atoi(id[2:])
	if err != nil {
		// this shouldn't happen though
		return StationID{}, err
	}
	return StationID{line: strings.ToUpper(id[:2]), number: number}, nil
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
		id, err := NewStationID(record[0])
		if err != nil {
			return nil, err
		}
		openingDate, err := time.Parse(openingDateFormat, record[2])
		if err != nil {
			return nil, err
		}
		final = append(final, Station{
			id:          id,
			name:        record[1],
			openingDate: openingDate,
		})
	}

	return final, nil
}
