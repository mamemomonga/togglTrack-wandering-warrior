package toggler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ta9mi141/toggl-go/track/toggl"
)

// "github.com/davecgh/go-spew/spew"

type Toggler struct {
	token          string
	workspaceId    int
	skipProjectIds []int
}

type NewOptions struct {
	Token          string
	WorkspaceId    int
	SkipProjectIds []int
}

func New(opt *NewOptions) *Toggler {
	t := &Toggler{}
	t.token = opt.Token
	t.workspaceId = opt.WorkspaceId
	t.skipProjectIds = opt.SkipProjectIds

	return t
}

func (t *Toggler) newAPIClient() *toggl.APIClient {
	return toggl.NewAPIClient(toggl.WithAPIToken(t.token))
}

type MonthlyEntries struct {
	Date   time.Time
	Exists bool
	Start  time.Time
	Stop   time.Time
	Rest   time.Duration
	Work   time.Duration
	Remain time.Duration
}

func (t *Toggler) Monthly(now time.Time, startStopRound time.Duration) []MonthlyEntries {

	monthlyEntries := []MonthlyEntries{}

	// 一ヶ月
	var startDateStr, endDateStr string
	var startDateTime, endDateTime time.Time
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}

	startDateTime = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, tokyo)
	endDateTime = startDateTime.AddDate(0, 1, 0)
	startDateStr = startDateTime.Format(time.DateOnly)
	endDateStr = endDateTime.Format(time.DateOnly)

	// 取得
	c := t.newAPIClient()
	timeEntries, err := c.GetTimeEntries(context.Background(), &toggl.GetTimeEntriesQuery{
		StartDate: &startDateStr,
		EndDate:   &endDateStr,
	})
	if err != nil {
		log.Fatal(err)
	}
	// spew.Dump(timeEntries)

	// Workspaceフィルタと日付ごとにまとめる
	var dailyEntries map[string][]*toggl.TimeEntry = map[string][]*toggl.TimeEntry{}
	for _, te := range timeEntries {
		if *te.WorkspaceID != t.workspaceId {
			continue
		}
		startDate := te.Start.Local().Format(time.DateOnly)
		dailyEntries[startDate] = append(dailyEntries[startDate], te)
	}

	// 日付ごとに処理する
	// fmt.Println(" 日付          | 開始  | 終了  | 休憩  | 作業  | 残稼働時間")
	for date := startDateTime; date.Before(endDateTime); date = date.AddDate(0, 0, 1) {
		dateStr := date.Format(time.DateOnly)

		// その日のデータがなければスキップ
		if _, ok := dailyEntries[dateStr]; !ok {
			// fmt.Printf("%s(%s)\n", dateStr, t.WeekdayJP(date))
			monthlyEntries = append(monthlyEntries, MonthlyEntries{
				Date:   date,
				Exists: false,
			})
			continue
		}

		// 項目の処理
		var timeStart time.Time
		var timeStop time.Time
		var workDuration time.Duration = 0

		for _, entry := range dailyEntries[dateStr] {
			// 休憩は無視
			skip := false
			for _, skipProjectId := range t.skipProjectIds {
				if entry.ProjectID != nil {
					if *entry.ProjectID == skipProjectId {
						skip = true
					}
				}
			}
			if skip {
				continue
			}

			// 稼働時間の処理
			duration, err := time.ParseDuration(fmt.Sprintf("%ds", *entry.Duration))
			if err != nil {
				log.Fatal(err)
			}
			duration = duration.Round(startStopRound)
			workDuration = workDuration + duration

			entryStart := *entry.Start
			entryStop := *entry.Stop
			entryStart = entryStart.Round(startStopRound).In(tokyo)
			entryStop = entryStop.Round(startStopRound).In(tokyo)

			// 休憩以外の場合
			if timeStart.IsZero() {
				timeStart = entryStart
			} else if entryStart.Before(timeStart) {
				timeStart = entryStart
			}
			if timeStop.IsZero() {
				timeStop = entryStop
			} else if entryStop.After(timeStop) {
				timeStop = entryStop
			}
		}
		// 開始から終了までの稼働時間
		totalDuration := timeStop.Sub(timeStart)
		// 休憩時間
		restDuration := totalDuration - workDuration

		monthlyEntries = append(monthlyEntries, MonthlyEntries{
			Date:   date,
			Exists: true,
			Start:  timeStart,
			Stop:   timeStop,
			Rest:   restDuration,
			Work:   workDuration,
		})
	}
	return monthlyEntries
}
