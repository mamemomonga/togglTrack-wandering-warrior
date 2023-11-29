package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/mamemomonga/togglTrack-wandering-warrior/src/configs"
	"github.com/mamemomonga/togglTrack-wandering-warrior/src/utils"
	"github.com/shiena/ansicolor"
	"github.com/wsxiaoys/terminal/color"
)

var (
	version  string
	revision string
)

var cfg *configs.Configs
var aw io.Writer

func main() {
	var err error
	aw = ansicolor.NewAnsiColorWriter(os.Stdout)

	var (
		configFile  = flag.String("config", "", "configファイル")
		monthOffset = flag.Int("month", 0, "今月からn月移動")
		demo        = flag.Bool("demo", false, "デモモード")
		offDates    = flag.Int("off", 0, "平日休暇日数")
	)
	flag.Parse()

	// 設定読込
	cfg, err = configs.New(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	version_string := ""
	if version != "" {
		version_string = fmt.Sprintf("v%s-%s", version, revision)
	}

	color.Fprintf(aw, "togglTrack-wandering-warrior %s\n", version_string)

	// 現在
	today := time.Now()
	if *demo {
		today = time.Date(2023, 10, 27, 0, 0, 0, 0, utils.TzTokyo())
		// today = time.Date(2023, 12, 27, 0, 0, 0, 0, utils.TzTokyo())
	}

	// 対象日
	target := today.AddDate(0, *monthOffset, 0)

	monthly(today, target, *offDates)

}
