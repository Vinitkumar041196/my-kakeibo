package utils

import (
	"time"
)

var loc *time.Location

func SetTimezone(tz string) error {
	location, err := time.LoadLocation(tz)
	if err != nil {
		return err
	}
	loc = location
	return nil
}

func TimeNow() time.Time {
	if loc == nil {
		loc = time.UTC
	}
	return time.Now().In(loc)
}
