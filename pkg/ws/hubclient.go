package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
)

type HubClient struct {
	Hub  *Hub
	Conn *websocket.Conn
	ID   int32
	send chan []byte
}

func NewHubClient(conn *websocket.Conn, hub *Hub) *HubClient {
	return &HubClient{
		Hub:  hub,
		Conn: conn,
		send: make(chan []byte, BufferSize),
		ID:   rand.Int31(),
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

func (c *HubClient) ReadPump() {
	defer func() {
		c.Hub.Unregister(*c)
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			log.Printf("Client %d disconnected", c.ID)
			break
		}

		log.Printf("From Client %d: %s", c.ID, message)

		c.Hub.Receive(message)
	}
}

func (c *HubClient) WritePump() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}

			w.Write(message)

			log.Printf("To Client %d: %s", c.ID, message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}
