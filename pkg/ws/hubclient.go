package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net"
)

type HubClient struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	id   int32

	// TODO: Need?
	remoteAddr net.Addr
}

func NewHubClient(conn *websocket.Conn, hub *Hub) *HubClient {
	return &HubClient{
		hub:  hub,
		conn: conn,
		// TODO: Set buffer size from constant
		send: make(chan []byte, 1024),
		id:   rand.Int31(),
	}
}

func (c *HubClient) Close() {
	close(c.send)
}

func (c *HubClient) Send(message []byte) {
	c.send <- message
}

func (c *HubClient) GetSendChannel() chan []byte {
	return c.send
}

func (c *HubClient) GetID() int32 {
	return c.id
}

func (c *HubClient) ReadPump() {
	defer func() {
		c.hub.Unregister(*c)
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.hub.Receive(message)
	}
}

func (c *HubClient) WritePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}

			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}
