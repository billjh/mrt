package main

import (
	"time"
)

// Station represents a MRT station in Singapore
type Station struct {
	code        string
	name        string
	openingDate time.Time
}

// MRTLine is the first two characters of Station.code
// For example, EW12 -> EW, and TE3 -> TE
// Empty string is returned when format is unexpected
func (s Station) MRTLine() string {
	if len(s.code) > 2 {
		return s.code[:2]
	}
	return ""
}
