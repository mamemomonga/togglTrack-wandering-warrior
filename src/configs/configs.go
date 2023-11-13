package configs

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"gopkg.in/yaml.v2"
)

const app_name = "togglTrack-wandering-warrior"

type Configs struct {
	Toggl          Toggl     `yaml:"toggl"`
	SkipProjectIds []int     `yaml:"skip_project_ids"`
	Worktimes      Worktimes `yaml:"worktimes"`
	Holidays       []DayOff  `yaml:"holidays"`
	DaysOff        []DayOff  `yaml:"daysoff"`
}

type Toggl struct {
	Token       string `yaml:"token"`
	WorkspaceId int    `yaml:"workspace_id"`
}

type DayOff struct {
	Date string `yaml:"date"`
	Name string `yaml:"name"`
}

type Worktimes struct {
	Max   time.Duration `yaml:"max"`
	Min   time.Duration `yaml:"min"`
	Start time.Time     `yaml:"start"`
	End   time.Time     `yaml:"end"`
	Rest  time.Duration `yaml:"rest"`
	Round time.Duration `yaml:"round"`
}

func (wt *Worktimes) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type alias struct {
		Max   string `yaml:"max"`
		Min   string `yaml:"min"`
		Start string `yaml:"start"`
		End   string `yaml:"end"`
		Rest  string `yaml:"rest"`
		Round string `yaml:"round"`
	}
	var err error
	var in alias
	if err := unmarshal(&in); err != nil {
		return err
	}

	wt.Start, err = time.Parse(time.TimeOnly, in.Start+":00")
	if err != nil {
		return fmt.Errorf("[config] worktimes.start エラー (%v)", err)
	}

	wt.End, err = time.Parse(time.TimeOnly, in.End+":00")
	if err != nil {
		return fmt.Errorf("[config] worktimes.end エラー (%v)", err)
	}

	wt.Max, err = time.ParseDuration(in.Max)
	if err != nil {
		return fmt.Errorf("[config] worktimes.maxエラー (%v)", err)
	}

	wt.Min, err = time.ParseDuration(in.Min)
	if err != nil {
		return fmt.Errorf("[config] worktimes.minエラー (%v)", err)
	}

	wt.Rest, err = time.ParseDuration(in.Rest)
	if err != nil {
		return fmt.Errorf("[config] worktimes.restエラー (%v)", err)
	}

	wt.Round, err = time.ParseDuration(in.Round)
	if err != nil {
		return fmt.Errorf("[config] worktimes.roundエラー (%v)", err)
	}
	return nil
}

func New(filename string) (t *Configs, err error) {

	if filename == "" {
		filename = path.Join(os.Getenv("HOME"), ".config", app_name, "config.yaml")
	}

	if !fileExists(filename) {
		createTemplateFile(filename)
		//		return nil, errors.New("error: configfile not exists")
	}

	t = &Configs{}
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(buf, t)
	if err != nil {
		return nil, err
	}
	// log.Printf("debug: [Read] %s", filename)
	if t.Toggl.Token == "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX" {
		return nil, fmt.Errorf("%s を修正してください", filename)
	}
	return t, nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func createTemplateFile(filename string) {
	hdoc := `
toggl:
# API Tokenは https://track.toggl.com/profile から取得できる
  token: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
# WorkspaceIDは「Manage Workspaces」から対象のWorkspaceを選び、URLの /workspaces/[数値]/ の[数値]の部分
  workspace_id: 0000000

skip_project_ids:
# 集計除外するプロジェクトIDを指定する
# ProjectIDは「MANAGE → Projects → プロジェクトをクリック」し、URLの /projects/[数値]/ の[数値]の部分
# 省略可
  - 0000000 # 休憩

worktimes:
  min: "140h"   # 最低稼働時間
  max: "180h"   # 最高稼働時間
  rest: "1h"    # 標準の休憩時間
  start: "9:00" # 標準の就業時間
  end: "18:00"  # 標準の終業時間
  round: "15m"  # 時間単位

# 日本の祝日
holidays:
  - { date: 2023-10-09, name: "スポーツの日" }
  - { date: 2023-11-03, name: "文化の日" }
  - { date: 2023-11-23, name: "勤労感謝の日" }
  - { date: 2023-12-31, name: "大晦日" }

# 独自に定義した休暇
daysoff:
  - { date: 2023-10-10 name: "バカンス" }
`
	_ = hdoc
	err := os.MkdirAll(path.Dir(filename), os.FileMode(0755))
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write([]byte(hdoc))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("---------------------------------\n")
	fmt.Printf("%s\n", filename)
	fmt.Printf("にテンプレートを書き出しました\n")
	fmt.Printf("このファイルを編集して再実行してください\n")
	fmt.Printf("---------------------------------\n")
	os.Exit(0)
}
