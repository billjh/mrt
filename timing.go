package main

import "time"

// Peak hours (6am-9am and 6pm-9pm on Mon-Fri)
// Night hours (10pm-6am on Mon-Sun)
// Non-Peak hours (all other times)
func isPeakHours(t time.Time) bool {
	if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
		return false
	}
	h := t.Hour()
	return h >= 6 && h < 9 || h >= 18 && h < 21
}

func isNightHours(t time.Time) bool {
	h := t.Hour()
	return h >= 22 || h < 6
}

// check if given MRT line stops operation at night hours
func stopAtNight(line string) bool {
	lines := []string{"DT", "CG", "CE"}
	for _, stoppedLine := range lines {
		if line == stoppedLine {
			return true
		}
	}
	return false
}
