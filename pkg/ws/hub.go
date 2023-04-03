package ws

import (
	"github.com/gorilla/websocket"
)

type RegisterHubClientEvent struct {
	Conn                   *websocket.Conn
	ClientRegistrationDone chan HubClient
}

type HubReceiveMessageEvent struct {
	Message interface{}
}

type HubSingleSendMessageEvent struct {
	Message  interface{}
	ClientID int32
}

type HubBroadcastMessageEvent struct {
	Message interface{}
}

type Hub struct {
	clients        map[int32]HubClient
	exist          map[string]HubClient
	register       chan RegisterHubClientEvent
	unregister     chan HubClient
	messageReceive chan HubReceiveMessageEvent
	singleSend     chan HubSingleSendMessageEvent
	broadcast      chan HubBroadcastMessageEvent
}

func NewHub() *Hub {
	return &Hub{
		clients:        make(map[int32]HubClient),
		exist:          make(map[string]HubClient),
		register:       make(chan RegisterHubClientEvent),
		unregister:     make(chan HubClient),
		messageReceive: make(chan HubReceiveMessageEvent),
		singleSend:     make(chan HubSingleSendMessageEvent),
		broadcast:      make(chan HubBroadcastMessageEvent),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case register := <-h.register:
			go h.handleRegisterHubClientEvent(&register)

		case hubClient := <-h.unregister:
			go h.handleUnregisterClientEvent(&hubClient)

		case message := <-h.singleSend:
			go h.handleSingleSendMessageEvent(message)

		case message := <-h.broadcast:
			go h.handleBroadcastMessageEvent(message)
		}
	}
}

func (h *Hub) Register(conn *websocket.Conn) HubClient {
	done := make(chan HubClient)

	h.register <- RegisterHubClientEvent{
		Conn:                   conn,
		ClientRegistrationDone: done,
	}

	return <-done
}

func (h *Hub) Unregister(c HubClient) {
	h.unregister <- c
}

func (h *Hub) Receive(message []byte) {
	h.messageReceive <- HubReceiveMessageEvent{
		Message: message,
	}
}

func (h *Hub) Send(clientID int32, message []byte) {
	h.singleSend <- HubSingleSendMessageEvent{
		Message:  message,
		ClientID: clientID,
	}
}

func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- HubBroadcastMessageEvent{
		Message: message,
	}
}

func (h *Hub) GetMessageReceiveChannel() <-chan HubReceiveMessageEvent {
	return h.messageReceive
}

func (h *Hub) GetRegisterChannel() <-chan RegisterHubClientEvent {
	return h.register
}

func (h *Hub) GetUnregisterChannel() <-chan HubClient {
	return h.unregister
}

func (h *Hub) ProcessSendMessages(c <-chan HubSingleSendMessageEvent) {
	for m := range c {
		go h.Send(m.ClientID, m.Message.([]byte))
	}
}

func (h *Hub) ProcessBroadcastMessages(c <-chan HubBroadcastMessageEvent) {
	for m := range c {
		go h.Broadcast(m.Message.([]byte))
	}
}

func (h *Hub) handleRegisterHubClientEvent(register *RegisterHubClientEvent) {
	client := NewHubClient(register.Conn, h)
	h.clients[client.ID] = *client
	register.ClientRegistrationDone <- *client
}

func (h *Hub) handleUnregisterClientEvent(c *HubClient) {
	if _, ok := h.clients[c.ID]; ok {
		delete(h.clients, c.ID)

		for k, v := range h.exist {
			if v == *c {
				delete(h.exist, k)
			}
		}

		c.Close()
	}
}

func (h *Hub) handleSingleSendMessageEvent(send HubSingleSendMessageEvent) {
	if client, ok := h.clients[send.ClientID]; ok {
		client.Send(send.Message.([]byte))
	}
}

func (h *Hub) handleBroadcastMessageEvent(broadcast HubBroadcastMessageEvent) {
	for _, client := range h.clients {
		client.Send(broadcast.Message.([]byte))
	}
}
