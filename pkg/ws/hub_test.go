package ws

import (
	"github.com/gorilla/websocket"
	"testing"
)

func TestHub_Register(t *testing.T) {
	hub := NewHub()

	go hub.Run()

	conn := &websocket.Conn{}
	client := hub.Register(conn)

	if client.GetID() == 0 {
		t.Errorf("Expected client ID to be non-zero, got %d", client.GetID())
	}

	if client.hub != hub {
		t.Errorf("Expected client hub to be the same as the argument, got %v", client.hub)
	}

	if client.conn != conn {
		t.Errorf("Expected client conn to be the same as the argument, got %v", client.conn)
	}
}
