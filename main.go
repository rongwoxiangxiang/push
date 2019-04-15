package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
	"webs/webservice"
	"webs/webservice/common"
)

const APP_NAME  = "webs"

func main()  {
	log.Printf("Push Application start [%s]", time.Now().String())

	webservice.InitConfig(filepath.Join(os.Getenv("GOPATH"), "src", APP_NAME, "application.json"))

	common.InitStats()

	webservice.InitConnMgr()

	webservice.InitWSServer()

	webservice.InitService()


}
