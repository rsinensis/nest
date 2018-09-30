package id

import (
	"log"
	"os"

	"github.com/rsinensis/nest/module/setting"
	"github.com/rsinensis/nest/util/snowflake"
)

var id *snowflake.Id

// InitId int id
func InitId() {
	cfg := setting.GetSetting()

	datacenter := cfg.Section("log").Key("Datacenter").MustInt64(1)
	worker := cfg.Section("log").Key("Worker").MustInt64(1)

	var err error
	id, err = snowflake.NewId(datacenter, worker, snowflake.GetIdTwepoch())
	if err != nil {
		log.Fatalf("snowflake NewId: %v", err)
		os.Exit(1)
	}
}

func GetId() *snowflake.Id {
	return id
}
