package webservice

import (
	"log"
	"time"
	"webs/webservice/common"
)

func Start(configs string) {
	if configs == "" {
		configs = "application.json"
	}
	log.Printf("Push Application start [%s]", time.Now().String())

	InitConfig(configs)

	common.InitStats()

	InitConnMgr()

	InitWSServer()

	InitService()
}
