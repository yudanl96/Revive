package chat

import (
	"github.com/gorilla/websocket"
)

// client is a single chat user
type client struct {
	socket  *websocket.Conn //web socket for this client
	receive chan []byte     //channel to receive message
	room    *room           // chat room
}

// read a message
func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

// write a message
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.receive {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
