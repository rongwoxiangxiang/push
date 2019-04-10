package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

type Client struct {
	id      string
	socket  *websocket.Conn
	send    chan []byte
	message *Message
}

func (c *Client) work(){
	go c.read()
	go c.write()
}

func (c *Client) read () {
	defer c.close()
	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Print(string(message))
		msg := Message{Sender:c.id}
		err = json.Unmarshal(message, &msg)
		if err != nil {
			fmt.Println("Unmarshal failed, ", err)
			return
		}
		if msg.Type == "" {
			msg.Type = MESSAGE_TYPE_TO_BROADCAST
		}
		manager.broadcast <- &msg
	}
}

func (c *Client) write() {
	defer c.close()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func (c *Client) close ()  {
	c.socket.Close()
}