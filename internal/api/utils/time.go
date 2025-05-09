package utils

import (
	"fmt"
	"time"
)

func TimeLocationPST() *time.Location {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		fmt.Println(err.Error())
	}

	if loc == nil {
		loc = time.UTC
	}
	return loc
}

func TimeIn(t time.Time, name string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err == nil {
		t = t.In(loc)
	}
	return t, err
}

func TimePST() time.Time {
	t, _ := TimeIn(time.Now(), "America/Los_Angeles")
	return t
}

func MonthDiff(a, b time.Time) (month int) {
	m := a.Month()
	for a.Before(b) {
		a = a.Add(time.Hour * 24 * 14)
		m2 := a.Month()
		if m2 != m {
			month++
		}
		m = m2
	}

	return
}

func GetStartDayOfWeek(tm time.Time, loc *time.Location) time.Time { //get monday 00:00:00
	weekday := time.Duration(tm.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	year, month, day := tm.Date()
	currentZeroDay := time.Date(year, month, day, 0, 0, 0, 0, loc)
	return currentZeroDay.Add(-1 * (weekday - 1) * 24 * time.Hour)
}
