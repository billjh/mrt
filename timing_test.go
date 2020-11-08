package main

import (
	"testing"
	"time"
)

// Peak hours (6am-9am and 6pm-9pm on Mon-Fri)
// Night hours (10pm-6am on Mon-Sun)
// Non-Peak hours (all other times)
func TestTravelPeriod(t *testing.T) {
	peak := "peak"
	night := "night"
	nonpeak := "nonpeak"

	timeLayout := "2006-01-02T15:04"
	sun := "2020-11-08T"
	mon := "2020-11-09T"
	fri := "2020-11-13T"
	sat := "2020-11-14T"
	for _, testCase := range []struct {
		t string
		h string
	}{
		{t: mon + "05:59", h: night},
		{t: fri + "05:59", h: night},
		{t: sat + "05:59", h: night},
		{t: sun + "05:59", h: night},

		{t: mon + "06:01", h: peak},
		{t: fri + "06:01", h: peak},
		{t: sat + "06:01", h: nonpeak},
		{t: sun + "06:01", h: nonpeak},

		{t: mon + "08:59", h: peak},
		{t: fri + "08:59", h: peak},
		{t: sat + "08:59", h: nonpeak},
		{t: sun + "08:59", h: nonpeak},

		{t: mon + "09:01", h: nonpeak},
		{t: fri + "09:01", h: nonpeak},
		{t: sat + "09:01", h: nonpeak},
		{t: sun + "09:01", h: nonpeak},

		{t: mon + "17:59", h: nonpeak},
		{t: fri + "17:59", h: nonpeak},
		{t: sat + "17:59", h: nonpeak},
		{t: sun + "17:59", h: nonpeak},

		{t: mon + "18:01", h: peak},
		{t: fri + "18:01", h: peak},
		{t: sat + "18:01", h: nonpeak},
		{t: sun + "18:01", h: nonpeak},

		{t: mon + "20:59", h: peak},
		{t: fri + "20:59", h: peak},
		{t: sat + "20:59", h: nonpeak},
		{t: sun + "20:59", h: nonpeak},

		{t: mon + "21:01", h: nonpeak},
		{t: fri + "21:01", h: nonpeak},
		{t: sat + "21:01", h: nonpeak},
		{t: sun + "21:01", h: nonpeak},

		{t: mon + "21:59", h: nonpeak},
		{t: fri + "21:59", h: nonpeak},
		{t: sat + "21:59", h: nonpeak},
		{t: sun + "21:59", h: nonpeak},

		{t: mon + "22:01", h: night},
		{t: fri + "22:01", h: night},
		{t: sat + "22:01", h: night},
		{t: sun + "22:01", h: night},
	} {
		p, err := time.Parse(timeLayout, testCase.t)
		if err != nil {
			t.Error(err)
		}
		var actual string
		switch {
		case isPeakHours(p):
			actual = peak
		case isNightHours(p):
			actual = night
		default:
			actual = nonpeak
		}
		if actual != testCase.h {
			t.Errorf("%s %s expected: %s, actual: %s", p.Weekday(), testCase.t, testCase.h, actual)
		}
	}
}
