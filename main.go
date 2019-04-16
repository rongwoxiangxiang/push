package main

import (
	"os"
	"path/filepath"
	"webs/webservice"
)

const APP_NAME  = "webs"

func main()  {
	webservice.Start(filepath.Join(os.Getenv("GOPATH"), "src", APP_NAME, "application.json"))
}
