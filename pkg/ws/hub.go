package ws

import (
	"github.com/gorilla/websocket"
)

type RegisterHubClientEvent struct {
	Conn                   *websocket.Conn
	ClientRegistrationDone chan HubClient
}

type HubMessageReceiveEvent struct {
	Message interface{}
}

type HubMessageSendEvent struct {
	Message  interface{}
	ClientID int32
}

type HubMessageBroadcastEvent struct {
	Message interface{}
}

type Hub struct {
	// Registered connections.
	clients        map[int32]HubClient
	exist          map[string]HubClient
	register       chan RegisterHubClientEvent
	unregister     chan HubClient
	messageReceive chan HubMessageReceiveEvent

	// TODO:
}

func NewHub() *Hub {
	return &Hub{
		clients:        make(map[int32]HubClient),
		exist:          make(map[string]HubClient),
		register:       make(chan RegisterHubClientEvent),
		unregister:     make(chan HubClient),
		messageReceive: make(chan HubMessageReceiveEvent),
	}
}

func (h *Hub) Register(c *websocket.Conn) HubClient {
	done := make(chan HubClient)

	h.register <- RegisterHubClientEvent{
		Conn:                   c,
		ClientRegistrationDone: done,
	}

	return <-done
}

func (h *Hub) Unregister(hubClient HubClient) {
	h.unregister <- hubClient
}

func (h *Hub) Run() {
	for {
		select {
		case register := <-h.register:
			go h.handleRegisterHubClientEvent(&register)

		case hubClient := <-h.unregister:
			for k, v := range h.exist {
				if v == hubClient {
					delete(h.exist, k)
				}
			}

			if _, ok := h.clients[hubClient.ID]; ok {
				delete(h.clients, hubClient.ID)
			}

			hubClient.Close()
		}

		// TODO: Process send and broadcast messages
	}
}

func (h *Hub) Receive(message []byte) {
	h.messageReceive <- HubMessageReceiveEvent{
		Message: message,
	}
}

func (h *Hub) Send(clientID int32, message []byte) {
	// TODO:
}

func (h *Hub) Broadcast(message []byte) {
	for _, client := range h.clients {
		client.Send(message)
	}
}

func (h *Hub) GetMessageReceiveChannel() <-chan HubMessageReceiveEvent {
	return h.messageReceive
}

func (h *Hub) GetRegisterChannel() <-chan RegisterHubClientEvent {
	return h.register
}

func (h *Hub) GetUnregisterChannel() <-chan HubClient {
	return h.unregister
}

func (h *Hub) ProcessSendMessages(c <-chan HubMessageSendEvent) {
	for m := range c {
		go h.Send(m.ClientID, m.Message.([]byte))
	}
}

func (h *Hub) ProcessBroadcastMessages(c <-chan HubMessageBroadcastEvent) {
	for m := range c {
		go h.Broadcast(m.Message.([]byte))
	}
}

func findHubClientByRemoteAddress(h *Hub, addr string) (HubClient, bool) {
	for _, client := range h.clients {
		if client.Conn.RemoteAddr().String() == addr {
			return client, true
		}
	}
	return HubClient{}, false
}

func (h *Hub) handleRegisterHubClientEvent(register *RegisterHubClientEvent) {
	client := NewHubClient(register.Conn, h)
	h.clients[client.ID] = *client
	register.ClientRegistrationDone <- *client
}
