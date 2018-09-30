package setting

import (
	"fmt"
	"log"
	"os"

	ini "gopkg.in/ini.v1"
)

var setting *ini.File

// InitSetting init setting
func InitSetting(mode string) {
	var err error
	cfgFile := fmt.Sprintf("conf/app_%v.ini", mode)
	setting, err = ini.Load(cfgFile)
	if err != nil {
		log.Fatalf("Fail to parse '%v': %v", cfgFile, err)
		os.Exit(1)
	}
}

// GetSetting return setting
func GetSetting() *ini.File {
	return setting
}
