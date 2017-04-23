package chat

import (
	"github.com/gorilla/websocket"
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
			return
		}
	}
}

func (c *client) writeLoop() {
	for msg := range c.recvc {
		if err := c.socket.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			return
		}
	}
}
