package chat

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/yvdinesh/go-blueprints/chat-server/trace"
	"log"
)

type Room struct {
	recvc   chan string
	joinc   chan *client
	removec chan *client
	clients map[*client]struct{}
	tracer  trace.Tracer
}

func NewRoom(options ...func(r *Room)) *Room {
	r := &Room{
		recvc:   make(chan string),
		clients: make(map[*client]struct{}),
		removec: make(chan *client),
		joinc:   make(chan *client),
	}
	for _, op := range options {
		op(r)
	}
	return r
}

func WithTracer(t trace.Tracer) func(r *Room) {
	return func(r *Room) {
		r.tracer = t
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
			r.tracer.Trace("New client joined")
		case c := <-r.removec:
			delete(r.clients, c)
			r.tracer.Trace("client left")
		case msg := <-r.recvc:
			r.tracer.Trace("Recieved message in the room: " + msg)
			for c := range r.clients {
				c.recvc <- msg
				r.tracer.Trace("-- sent to client")
			}
		}
	}
}
