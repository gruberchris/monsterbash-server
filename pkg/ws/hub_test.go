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

	time.Sleep(10 * time.Millisecond)

	hub.clientsMu.RLock()

	if _, ok := hub.clients[client.ID]; ok {
		t.Errorf("Expected client to be removed from the hub, but it still exists")
	}

	defer hub.clientsMu.RUnlock()
}

func TestHub_Receive(t *testing.T) {
	hub := NewHub()

	go hub.Run()

	conn := &websocket.Conn{}
	hub.Register(conn)

	hub.Receive([]byte("test"))

	time.Sleep(10 * time.Millisecond)

	select {
	case message := <-hub.messageReceive:
		messageBytes := message.Message.([]byte)
		if string(messageBytes) != "test" {
			t.Errorf("Expected message to be 'test', got %v", message.Message)
		}
	default:
		t.Errorf("Expected message to be received")
	}
}

func TestHub_Send(t *testing.T) {
	hub := NewHub()

	go hub.Run()

	conn := &websocket.Conn{}
	client := hub.Register(conn)

	hub.Send(client.ID, []byte("test"))

	time.Sleep(10 * time.Millisecond)

	select {
	case message := <-client.send:
		if string(message) != "test" {
			t.Errorf("Expected message to be 'test', got %v", message)
		}
	default:
		t.Errorf("Expected message to be received")
	}
}

func TestHub_Broadcast(t *testing.T) {
	hub := NewHub()

	go hub.Run()

	conn := &websocket.Conn{}
	client := hub.Register(conn)

	conn2 := &websocket.Conn{}
	client2 := hub.Register(conn2)

	hub.Broadcast([]byte("test"))

	time.Sleep(10 * time.Millisecond)

	message := <-client.send

	if string(message) != "test" {
		t.Errorf("Expected Client One message to be 'test', got %v", message)
	}

	message = nil

	message = <-client2.send

	if string(message) != "test" {
		t.Errorf("Expected Client Two message to be 'test', got %v", message)
	}
}

func TestHub_ProcessSendMessages(t *testing.T) {
	hub := NewHub()

	go hub.Run()

	conn := &websocket.Conn{}
	client := hub.Register(conn)

	singleSend := make(chan HubSingleSendMessageEvent, BufferSize)

	singleSend <- HubSingleSendMessageEvent{
		Message:  []byte("test"),
		ClientID: client.ID,
	}

	// Close the channel to avoid a deadlock
	close(singleSend)

	hub.ProcessSendMessages(singleSend)

	time.Sleep(10 * time.Millisecond)

	select {
	case message := <-client.send:
		if string(message) != "test" {
			t.Errorf("Expected message to be 'test', got %v", message)
		}
	}
}

func TestHub_ProcessBroadcastMessages(t *testing.T) {
	hub := NewHub()

	go hub.Run()

	conn := &websocket.Conn{}
	client := hub.Register(conn)

	conn2 := &websocket.Conn{}
	client2 := hub.Register(conn2)

	broadcastSend := make(chan HubBroadcastMessageEvent, BufferSize)

	broadcastSend <- HubBroadcastMessageEvent{
		Message: []byte("test"),
	}

	// Close the channel to avoid a deadlock
	close(broadcastSend)

	hub.ProcessBroadcastMessages(broadcastSend)

	time.Sleep(10 * time.Millisecond)

	message := <-client.send

	if string(message) != "test" {
		t.Errorf("Expected Client One message to be 'test', got %v", message)
	}

	message = nil

	message = <-client2.send

	if string(message) != "test" {
		t.Errorf("Expected Client Two message to be 'test', got %v", message)
	}
}
