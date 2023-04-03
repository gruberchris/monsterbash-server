package ws

import (
	"github.com/gorilla/websocket"
	"testing"
	"time"
)

func TestHub_Register(t *testing.T) {
	hub := NewHub()

	go hub.Run()

	conn := &websocket.Conn{}
	client := hub.Register(conn)

	if client.ID == 0 {
		t.Errorf("Expected client ID to be non-zero, got %d", client.ID)
	}

	if client.Hub != hub {
		t.Errorf("Expected client hub to be the same as the argument, got %v", client.Hub)
	}

	if client.Conn != conn {
		t.Errorf("Expected client conn to be the same as the argument, got %v", client.Conn)
	}
}

func TestHub_Unregister(t *testing.T) {
	hub := NewHub()

	go hub.Run()

	conn := &websocket.Conn{}
	client := hub.Register(conn)

	hub.Unregister(client)

	time.Sleep(1 * time.Second)

	hub.clientsMu.RLock()

	if _, ok := hub.clients[client.ID]; ok {
		t.Errorf("Expected client to be removed from the hub, but it still exists")
	}

	defer hub.clientsMu.RUnlock()
}
