package main

import (
	"fmt"
	"time"

	//	"github.com/davecgh/go-spew/spew"

	"github.com/mamemomonga/togglTrack-wandering-warrior/src/utils"
)

type MonthlyCalendarDay struct {
	Text      string
	Color     string
	OffName   string
	IsHoliday bool
	IsDayOff  bool
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

	type kvht struct {
		text    string
		holiday bool
	}

	kvh := make(map[string]kvht, len(cfg.Holidays)+len(cfg.DaysOff))
	for _, v := range cfg.Holidays {
		kvh[v.Date] = kvht{text: v.Name, holiday: true}
	}
	for _, v := range cfg.DaysOff {
		kvh[v.Date] = kvht{text: v.Name, holiday: false}
	}

	for date := start; date.Before(end); date = date.AddDate(0, 0, 1) {

		ca := MonthlyCalendarDay{
			IsDayOff:  false,
			IsHoliday: false,
			IsWeekday: false,
		}

		week := date.Weekday()
		weekJp := []string{"日", "月", "火", "水", "木", "金", "土"}
		dateStr := date.Format(time.DateOnly)
		ca.Text = fmt.Sprintf("%s(%s)", dateStr, weekJp[week])

		if v, ok := kvh[dateStr]; ok {
			ca.IsDayOff = true
			ca.IsWeekday = false
			ca.IsHoliday = v.holiday
			ca.OffName = v.text
			if v.holiday {
				ca.Color = "red"
			} else {
				ca.Color = "blue"
			}
			appendCal(ca, dateStr)
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
