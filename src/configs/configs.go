package configs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Configs struct {
	Toggl          Toggl     `yaml:"toggl"`
	SkipProjectIds []int     `yaml:"skip_project_ids"`
	Holidays       []Holiday `yaml:"holidays"`
	Worktimes      Worktimes `yaml:"worktimes"`
}

type Toggl struct {
	Token       string `yaml:"token"`
	WorkspaceId int    `yaml:"workspace_id"`
}

type Holiday struct {
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
	if !fileExists(filename) {
		return nil, errors.New("error: configfile not exists")
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
	log.Printf("debug: [Read] %s", filename)
	return t, nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
