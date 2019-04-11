package ws

import (
	"fmt"
	"log"
	"net/http"
	"github.com/satori/go.uuid"
	"github.com/gorilla/websocket"
)

type ClientManager struct {
	clients    map[string]*Client
	message    chan *Message
	register   chan *Client
	unregister chan *Client
}

func (manager *ClientManager) CreateClient(res http.ResponseWriter, req *http.Request) (client *Client, err error) {
	var conn *websocket.Conn
	conn, err = (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client = &Client{id: uuid.Must(uuid.NewV4()).String(), socket: conn, send: make(chan []byte)}
	client.work()
	manager.register <- client
	return
}

func (manager *ClientManager) Start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn.id] = conn
			fmt.Println("Conn: new user connected: \n%v", conn)
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn.id]; ok {
				manager.disconnected(conn)
			}
		case message := <-manager.message:
			if message.MessageHandler == nil {
				message.DefaultMessageHandler(manager)
				continue
			}
			message.MessageHandler(manager)
		}
	}
}

func (manager *ClientManager) disconnected(conn *Client) {
	close(conn.send)
	delete(manager.clients, conn.id)
	log.Println("Conn: one user disconnected: \n%v", conn)
}

