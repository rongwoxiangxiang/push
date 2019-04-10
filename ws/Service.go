package ws

import (
	"fmt"
	"net/http"
	"strings"
)

type Service struct {
	host string
	port string
}

var manager = &ClientManager{
	broadcast:  make(chan *Message),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[string]*Client),
}

func (this *Service) SetServer (host, port string) *Service {
	this.host = host
	this.port = port
	return this
}

func (this *Service) DefaultServer () {
	this.host = "0.0.0.0"
	this.port = "12345"
}

func WsDefaultHandler(response http.ResponseWriter, request *http.Request) {
	client, err := manager.CreateClient(response, request)
	if err != nil {
		fmt.Println(err)
	}
	if client != nil {

	}
}

func (this *Service) Server(pattern string) {
	if this.host == "" {
		this.DefaultServer()
	}
	if !strings.HasPrefix(pattern, "/") {
		pattern = "/" + pattern
	}
	go manager.Start()
	http.HandleFunc(pattern, WsDefaultHandler)
	http.ListenAndServe(this.host + ":" + this.port, nil)
}