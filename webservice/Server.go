package webservice

import (
	"log"
	"time"
)

func Run() {
	log.Printf("Push Application start [%s]", time.Now().String())

	InitConnMgr()

	InitWSServer()

	InitService()
}
