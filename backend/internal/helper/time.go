package helper

import "time"

func StartOfPreviousMonth(t time.Time) time.Time {
	prev := t.AddDate(0, -1, 0)

	return time.Date(
		prev.Year(),
		prev.Month(),
		1,
		0, 0, 0, 0,
		t.Location(),
	)
}
