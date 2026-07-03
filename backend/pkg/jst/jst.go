package jst

import "time"

var location *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic("failed to load Asia/Tokyo timezone: " + err.Error())
	}
	location = loc
}

// Location returns the Asia/Tokyo timezone location.
func Location() *time.Location {
	return location
}

// Now returns the current time in JST.
func Now() time.Time {
	return time.Now().In(location)
}

// TodayString returns today's date as YYYY-MM-DD in JST.
func TodayString() string {
	return Now().Format("2006-01-02")
}
