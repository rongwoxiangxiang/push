package main

import (
	"os"
	"path/filepath"
	"webs/webservice"
	"webs/webservice/common"
)

const APP_NAME  = "webs"

func main()  {

	webservice.InitConfig(filepath.Join(os.Getenv("GOPATH"), "src", APP_NAME, "application.json"))

	common.InitStats()

	webservice.InitConnMgr()

	webservice.InitWSServer()

	webservice.InitService()
}
