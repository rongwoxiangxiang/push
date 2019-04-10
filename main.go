package main

import (
	"webs/ws"
)

func main()  {
	(&ws.Service{}).Server("/test")
}
