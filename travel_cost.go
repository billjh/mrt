package main

import (
	"time"
)

// TravelCost is an interface type for getting cost of travel between Stations
type TravelCost interface {
	Interchange() Weight
	OnLine(line string) Weight
}

// TravelCostByStop gives cost 1 for both interchange and travel on line
type TravelCostByStop struct{}

// Interchange implements TravelCost interface
func (c TravelCostByStop) Interchange() Weight { return 1 }

// OnLine implements TravelCost interface
func (c TravelCostByStop) OnLine(_ string) Weight { return 1 }

// TravelCostByTime contains costs for interchange and travel on line
type TravelCostByTime struct {
	interchange Weight
	lines       map[string]Weight
	lineDefault Weight
}

// Interchange implements TravelCost interface
func (c TravelCostByTime) Interchange() Weight {
	return c.interchange
}

// OnLine implements TravelCost interface
func (c TravelCostByTime) OnLine(line string) Weight {
	if w, ok := c.lines[line]; ok {
		return w
	}
	return c.lineDefault
}

var travelCostPeakHours = TravelCostByTime{
	interchange: 15,
	lines: map[string]Weight{
		"NS": 12,
		"NE": 12,
	},
	lineDefault: 10,
}

var travelCostNightHours = TravelCostByTime{
	interchange: 10,
	lines: map[string]Weight{
		"TE": 8,
	},
	lineDefault: 10,
}

var travelCostNonPeakHours = TravelCostByTime{
	interchange: 10,
	lines: map[string]Weight{
		"DT": 8,
		"TE": 8,
	},
	lineDefault: 10,
}

// getTravelCostByTime is a helper function to get travel cost based on time period
func getTravelCostByTime(t time.Time) TravelCost {
	switch {
	case isPeakHours(t):
		return travelCostPeakHours
	case isNightHours(t):
		return travelCostNightHours
	default:
		return travelCostNonPeakHours
	}
}
