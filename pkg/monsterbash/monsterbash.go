package monsterbash

import (
	"log"
	"monsterbash-server/pkg/ws"
	"time"
)

var BufferSize = 1024

type MonsterBash struct {
	ticker                  *time.Ticker
	quitChannel             chan bool
	sendMessageChannel      chan ws.HubSingleSendMessageEvent
	broadcastMessageChannel chan ws.HubBroadcastMessageEvent
}

func NewMonsterBash() *MonsterBash {
	return &MonsterBash{
		ticker:                  time.NewTicker(1 * time.Second),
		quitChannel:             make(chan bool),
		sendMessageChannel:      make(chan ws.HubSingleSendMessageEvent, BufferSize),
		broadcastMessageChannel: make(chan ws.HubBroadcastMessageEvent, BufferSize),
	}
}

func (mb *MonsterBash) Run() {
	for {
		select {
		case <-mb.ticker.C:
			mb.Update()

		case <-mb.quitChannel:
			mb.ticker.Stop()
			return
		}
	}
}

func (mb *MonsterBash) Update() {
	// TODO: Update the game state. Loop through all arenas and update state in them
}

func (mb *MonsterBash) ProcessInput(c <-chan ws.HubReceiveMessageEvent) {
	for m := range c {
		switch m.Message.(type) {
		// TODO: Handle game messages from the player
		}
	}
}

func (mb *MonsterBash) ProcessUnregisteredPlayers(c <-chan ws.HubClient) {
	for hubClient := range c {
		go mb.removePlayer(hubClient)
	}
}

func (mb *MonsterBash) ProcessNewPlayers(c <-chan ws.PlayerJoinGameEvent) {
	for m := range c {
		hubClient := m.PlayerHubClient
		go mb.connectPlayer(hubClient)
	}
}

func (mb *MonsterBash) GetQuitChannel() <-chan bool {
	return mb.quitChannel
}

func (mb *MonsterBash) GetSendMessageChannel() <-chan ws.HubSingleSendMessageEvent {
	return mb.sendMessageChannel
}

func (mb *MonsterBash) GetBroadcastMessageChannel() <-chan ws.HubBroadcastMessageEvent {
	return mb.broadcastMessageChannel
}

func (mb *MonsterBash) connectPlayer(client *ws.HubClient) {
	log.Printf("New player %d joined", client.ID)

	go client.WritePump()
	go client.ReadPump()

	// TODO: Send player their list of available arenas
}

func (mb *MonsterBash) removePlayer(client ws.HubClient) {
	// TODO: Remove the player from the game

	// TODO: Clean up any active arenas the player was in

	// TODO: Send player a message that they have been removed from the game
}
