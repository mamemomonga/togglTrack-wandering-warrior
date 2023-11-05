package main

import (
	"flag"
	"io"
	"log"
	"os"
	"time"

	"example.com/toggl/configs"
	"github.com/shiena/ansicolor"
)

var cfg *configs.Configs
var aw io.Writer

func main() {

	var err error

	var (
		configFile  = flag.String("config", "./config.yaml", "configファイル")
		monthOffset = flag.Int("month", 0, "今月からn月戻る")
	)
	flag.Parse()

	// 設定読込
	cfg, err = configs.New(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	// 対象日
	targetDate := time.Now().AddDate(0, -*monthOffset, 0)

	monthOffsetMode := false
	if *monthOffset != 0 {
		monthOffsetMode = true
	}

	aw = ansicolor.NewAnsiColorWriter(os.Stdout)

	monthly(targetDate, monthOffsetMode)

}