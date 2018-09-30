package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rsinensis/nest/module/logger"
	"github.com/rsinensis/nest/module/setting"
	"github.com/rsinensis/nest/module/id"
)

var (
	// BuildVersion from git tag
	BuildVersion string
	// BuildTime from make time
	BuildTime string
	// BuildMode from make mode
	BuildMode string
)

// dispay app info
func info() {
	var v bool
	flag.BoolVar(&v, "v", false, "show version and exit")
	flag.Parse()

	if v {
		log.Println(fmt.Sprintf("\nBuildVersion: %v\n   BuildTime: %v\n   BuildMode: %v", BuildVersion, BuildTime, BuildMode))
		os.Exit(0)
	}

	if len(BuildMode) == 0 {
		BuildMode = "dev"
	}
}

func main() {
	info()

	setting.InitSetting(BuildMode)
	logger.InitLogger(BuildMode)
	id.InitId()
}
