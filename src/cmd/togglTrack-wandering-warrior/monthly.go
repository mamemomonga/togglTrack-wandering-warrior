package main

import (
	"time"

	"github.com/mamemomonga/togglTrack-wandering-warrior/src/toggler"
	"github.com/mamemomonga/togglTrack-wandering-warrior/src/utils"
	"github.com/wsxiaoys/terminal/color"
)

func monthly(today, target time.Time) {

	weekdayTotal := 0
	weekdayRemain := 0

	tgr := toggler.New(
		&toggler.NewOptions{
			Token:          cfg.Toggl.Token,
			WorkspaceId:    cfg.Toggl.WorkspaceId,
			SkipProjectIds: cfg.SkipProjectIds,
		})

	workTimeTotal := time.Duration(0)
	monthly := tgr.Monthly(target, cfg.Worktimes.Round)

	color.Fprintf(aw,
		"@{bW} %s         @{|}|@{bW} %s  @{|}|@{bW} %s  @{|}|@{bW} %s   @{|}|@{bW} %s   @{|}|@{bW} %s  @{|}|@{bW} %s @{|}\n",
		" 日付", "開始", "終了", "休憩", "稼働", "総稼働", "最低達成",
	)

	for daynum, cal := range monthlyCalendar(today, target) {
		ent := monthly[daynum]
		switch cal.Color {
		case "red":
			color.Fprintf(aw, "@{wR}%s@{|}", cal.Text)
		case "blue":
			color.Fprintf(aw, "@{wB}%s@{|}", cal.Text)
		case "yellow":
			color.Fprintf(aw, "@{bY}%s@{|}", cal.Text)
		default:
			color.Fprintf(aw, "%s@{|}", cal.Text)
		}

		if cal.IsWeekday {
			weekdayTotal++
			if ent.Date.After(today) {
				weekdayRemain++
			}
		}

		if ent.Exists {

			workTimeTotal = workTimeTotal + ent.Work
			color.Fprintf(aw, " | %5s | %5s | %6s | %6s | %7s | %5.1f%%\n",
				ent.Start.Format(time.TimeOnly)[0:5],
				ent.Stop.Format(time.TimeOnly)[0:5],
				utils.DurationHourMin(ent.Rest),
				utils.DurationHourMin(ent.Work),
				utils.DurationHourMin(workTimeTotal),
				(workTimeTotal.Seconds()/cfg.Worktimes.Min.Seconds())*100,
			)
		} else {
			color.Fprintf(aw, " %s\n", cal.OffName)
		}
	}

	monthlySummary(today, target, weekdayTotal, weekdayRemain, workTimeTotal)
}

func monthlySummary(now, target time.Time, weekdayTotal, weekdayRemain int, workTimeTotal time.Duration) {
	color.Fprintf(aw, "@{bW}       %04d年%02d月の概要       @{|}\n", target.Year(), target.Month())

	color.Fprintf(aw, "今日は@{!}%04d年%02d月%02d日@{|}です\n", now.Year(), now.Month(), now.Day())
	color.Fprintf(aw, "@{!}%02d@{|}月の平日は@{!}%d日間@{|}で@{!}残り%d日@{|}で", target.Month(), weekdayTotal, weekdayRemain)
	color.Fprintf(aw, "@{!}%5.2f%%@{|}が経過します\n", (float64(weekdayTotal-weekdayRemain)/float64(weekdayTotal))*100)
	color.Fprintf(aw, "現在の総稼働時間は@{!}%s@{|}です\n", utils.DurationHourMin(workTimeTotal))

	guessRemain := cfg.Worktimes.End.Sub(cfg.Worktimes.Start).Hours() * float64(weekdayRemain)
	guessTotal := guessRemain + workTimeTotal.Hours()

	color.Fprintf(aw, "残り@{!}%d日間@{|}を@{!}%s~%s(休憩%.0f時間)@{|}で稼働した場合\n",
		weekdayRemain,
		cfg.Worktimes.Start.Format(time.TimeOnly)[:5],
		cfg.Worktimes.End.Format(time.TimeOnly)[:5],
		cfg.Worktimes.Rest.Hours(),
	)
	color.Fprintf(aw, "  残総稼働時間は@{!}%.0f時間@{|}となり\n",
		guessRemain,
	)

	color.Fprintf(aw, "    最低稼働@{!}%.0f時間@{|}の@{!}%6.2f%%@{|}\n",
		cfg.Worktimes.Min.Hours(),
		(guessTotal/cfg.Worktimes.Min.Hours())*100,
	)
	color.Fprintf(aw, "    最高稼働@{!}%.0f時間@{|}の@{!}%6.2f%%@{|}\n",
		cfg.Worktimes.Max.Hours(),
		(guessTotal/cfg.Worktimes.Max.Hours())*100,
	)

	if guessTotal < cfg.Worktimes.Min.Hours() {
		color.Fprintf(aw, "    @{bY} 未達が予測されます @{|}\n")
	} else if guessTotal > cfg.Worktimes.Max.Hours() {
		color.Fprintf(aw, "    @{wR} 超過が予測されます @{|}\n")
	} else {
		color.Fprintf(aw, "    @{wB} 達成が予測されます @{|}\n")
	}

	color.Fprintf(aw, "今週も勤労に勤しみましょう。\n")
}