package tasks

import (
	"time"
)

const TIME_FORMAT = "2006-01-02T15:04:05.000+0000"

func convertUTCString(t string) time.Time {
	tm, err := time.Parse(TIME_FORMAT, t)
	if err != nil {
		return time.Time{}
	}
	return tm.Local()
}

func convertLocalTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(TIME_FORMAT)
}
