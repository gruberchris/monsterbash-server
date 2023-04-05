package ws

import (
	"github.com/gorilla/websocket"
	"sync"
)

var BufferSize = 1024

type RegisterHubClientEvent struct {
	Conn                   *websocket.Conn
	ClientRegistrationDone chan HubClient
}

type PlayerJoinGameEvent struct {
	PlayerHubClient *HubClient
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
	clientsMu      sync.RWMutex
	clients        map[int32]HubClient
	existMu        sync.RWMutex
	exist          map[string]HubClient
	register       chan RegisterHubClientEvent
	playerJoinGame chan PlayerJoinGameEvent
	unregister     chan HubClient
	messageReceive chan HubReceiveMessageEvent
	singleSend     chan HubSingleSendMessageEvent
	broadcast      chan HubBroadcastMessageEvent
}

func NewHub() *Hub {
	return &Hub{
		clients:        make(map[int32]HubClient),
		exist:          make(map[string]HubClient),
		register:       make(chan RegisterHubClientEvent, BufferSize),
		playerJoinGame: make(chan PlayerJoinGameEvent, BufferSize),
		unregister:     make(chan HubClient, BufferSize),
		messageReceive: make(chan HubReceiveMessageEvent, BufferSize),
		singleSend:     make(chan HubSingleSendMessageEvent, BufferSize),
		broadcast:      make(chan HubBroadcastMessageEvent, BufferSize),
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

func (h *Hub) GetPlayerJoinGameChannel() <-chan PlayerJoinGameEvent {
	return h.playerJoinGame
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

	h.clientsMu.Lock()
	h.clients[client.ID] = *client
	defer h.clientsMu.Unlock()

	h.playerJoinGame <- PlayerJoinGameEvent{
		PlayerHubClient: client,
	}

	register.ClientRegistrationDone <- *client
}

func (h *Hub) handleUnregisterClientEvent(c *HubClient) {
	if _, ok := h.clients[c.ID]; ok {
		h.clientsMu.Lock()
		delete(h.clients, c.ID)
		defer h.clientsMu.Unlock()

		h.existMu.Lock()

		for k, v := range h.exist {
			if v == *c {
				delete(h.exist, k)
			}
		}

		defer h.existMu.Unlock()

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
