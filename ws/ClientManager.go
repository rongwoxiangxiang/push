package ws

import (
	"fmt"
	"net/http"
	"github.com/satori/go.uuid"
	"github.com/gorilla/websocket"
)

type ClientManager struct {
	clients    map[string]*Client
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
}

func (manager *ClientManager) CreateClient(res http.ResponseWriter, req *http.Request) (client *Client, err error) {
	var conn *websocket.Conn
	conn, err = (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		fmt.Println(err)
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
		case message := <-manager.broadcast:
			switch message.Type {
			case MESSAGE_TYPE_TO_SINGLE_USER:
				conn, ok := manager.clients[message.Recipient]
				if !ok {
					manager.disconnected(conn)
				}
				conn.send <- []byte(message.Content)
			case MESSAGE_TYPE_TO_BROADCAST:
				fmt.Println(111111)
				for _, conn := range manager.clients {
					select {
					case conn.send <- []byte(message.Content):
					default:
						manager.disconnected(conn)
					}
				}
			}
		}
	}
}

func (manager *ClientManager) send(message []byte, ignore *Client) {
	for _, conn := range manager.clients {
		if conn != ignore {
			conn.send <- message
		}
	}
}

func (manager *ClientManager) disconnected(conn *Client) {
	close(conn.send)
	delete(manager.clients, conn.id)
	fmt.Println("Conn: one user disconnected: \n%v", conn)
}

