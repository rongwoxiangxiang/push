package ws

import (
	"log"
)

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"to,omitempty"`
	Content   string `json:"content,omitempty"`
	Type	  string `json:"type,omitempty"`
	MessageHandler func (manager *ClientManager)
}

const (
	MESSAGE_TYPE_TO_SINGLE_USER = "single"	//单发
	MESSAGE_TYPE_TO_GROUP_USER  = "group"	//组发
	MESSAGE_TYPE_TO_MANY_USER   = "many" 	//群发
	MESSAGE_TYPE_TO_BROADCAST   = "broadcast" //广播
)

func (message *Message) DefaultMessageHandler(manager *ClientManager) {
	switch message.Type {
	case MESSAGE_TYPE_TO_SINGLE_USER:
		conn, ok := manager.clients[message.Recipient]
		if !ok {
			log.Println("Conn: message to single err: \n%v", message)
			return
		}
		conn.send <- []byte(message.Content)
	case MESSAGE_TYPE_TO_BROADCAST:
		for _, conn := range manager.clients {
			select {
			case conn.send <- []byte(message.Content):
			default:
				manager.disconnected(conn)
			}
		}
	}
}

func (message *Message) SetMessageHandler (handler func(manager *ClientManager)()) {
	message.MessageHandler = handler
}