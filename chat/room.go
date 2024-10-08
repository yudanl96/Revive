package chat

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	clients map[*client]bool //clients in this room
	join    chan *client     // channel for client wishing to joing the room
	leave   chan *client     // channel for lcient wishing to leave the room
	forward chan []byte      //channle that holds incomming messages, will be forwarded to other client
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func NewRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) Run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.receive)
		case msg := <-r.forward:
			for client := range r.clients {
				client.receive <- msg
			}
		}
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	socket, err := upgrader.Upgrade(w, request, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	client := &client{
		socket:  socket,
		receive: make(chan []byte, messageBufferSize),
		room:    r,
	}

	r.join <- client
	defer func() {
		r.leave <- client
	}()
	go client.write()
	client.read()

}
