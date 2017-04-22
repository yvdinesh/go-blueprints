package chat

import (
	"net/http"

	"github.com/gorilla/websocket"
	"log"
)

type Room struct {
	recvc   chan string
	joinc   chan *client
	removec chan *client
	clients map[*client]struct{}
}

func NewRoom() *Room {
	return &Room{
		recvc:   make(chan string),
		clients: make(map[*client]struct{}),
		removec: make(chan *client),
		joinc:   make(chan *client),
	}
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, w.Header())
	if err != nil {
		log.Fatalf("ServeHttp: %v", err.Error())
	}
	c := &client{
		recvc:  make(chan string),
		room:   r,
		socket: socket,
	}
	r.joinc <- c
	defer func() { r.removec <- c }()
	go c.writeLoop()
	c.readLoop()
}

func (r *Room) Run() {
	for {
		select {
		case c := <-r.joinc:
			r.clients[c] = struct{}{}
		case c := <-r.removec:
			delete(r.clients, c)
		case msg := <-r.recvc:
			for c := range r.clients {
				c.recvc <- msg
			}
		}
	}
}
