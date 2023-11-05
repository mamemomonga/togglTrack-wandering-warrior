package utils

import (
	"log"
	"time"
)

func TzTokyo() *time.Location {
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}
	return tokyo
}

func MonthRange(target time.Time) (start, end time.Time) {
	start = time.Date(target.Year(), target.Month(), 1, 0, 0, 0, 0, TzTokyo())
	end = start.AddDate(0, 1, 0)
	return start, end
}

func DurationHourMin(dur time.Duration) string {
	tw := dur.Round(time.Minute).String()
	buf := tw[:(len(tw) - 2)]
	if buf == "" {
		buf = "0s"
	}
	return buf
}
func DurationClockFmtHourMin(dur time.Duration) string {
	return time.Unix(0, 0).UTC().Add(time.Duration(dur)).Format(time.TimeOnly)[0:5]
}
