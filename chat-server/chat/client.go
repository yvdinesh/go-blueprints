package chat

import (
	"github.com/gorilla/websocket"
	"log"
)

type client struct {
	socket *websocket.Conn
	recvc  chan string
	room   *Room
}

func (c *client) readLoop() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.recvc <- string(msg)
		} else {
			log.Fatalf("Error while reading from socket: %v", err.Error())
		}
	}
}

func (c *client) writeLoop() {
	for msg := range c.recvc {
		if err := c.socket.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			log.Fatalf("Error while sending to socket: %v\n", err.Error())
		}
	}
}
