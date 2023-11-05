package main

import (
	"fmt"
	"time"

	"github.com/mamemomonga/togglTrack-wandering-warrior/src/utils"
)

type MonthlyCalendarDay struct {
	Text      string
	Color     string
	OffName   string
	IsHoliday bool
	IsWeekday bool
}

func monthlyCalendar(today, target time.Time) []MonthlyCalendarDay {
	cal := []MonthlyCalendarDay{}

	appendCal := func(ca MonthlyCalendarDay, dateStr string) {
		if today.Format(time.DateOnly) == dateStr {
			ca.Color = "yellow"
		}
		cal = append(cal, ca)
	}

	start, end := utils.MonthRange(target)

	for date := start; date.Before(end); date = date.AddDate(0, 0, 1) {

		ca := MonthlyCalendarDay{
			IsHoliday: false,
			IsWeekday: false,
		}

		week := date.Weekday()
		weekJp := []string{"日", "月", "火", "水", "木", "金", "土"}
		dateStr := date.Format(time.DateOnly)
		ca.Text = fmt.Sprintf("%s(%s)", dateStr, weekJp[week])

		for _, ho := range cfg.Holidays {
			if dateStr == ho.Date {
				ca.IsHoliday = true
				ca.IsWeekday = false
				ca.OffName = ho.Name
				ca.Color = "red"
				appendCal(ca, dateStr)
				continue
			}
		}
		if ca.IsHoliday {
			continue
		}
		switch week {
		case 0: // 日曜日
			ca.OffName = "休日"
			ca.Color = "red"
		case 6: // 土曜日
			ca.OffName = "休日"
			ca.Color = "blue"
		default: // 平日
			ca.IsWeekday = true
		}
		appendCal(ca, dateStr)
	}

	return cal

}
