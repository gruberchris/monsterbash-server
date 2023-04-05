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
	//TODO:
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
	// TODO: Update the game state
}

func (mb *MonsterBash) ProcessInput(c <-chan ws.HubReceiveMessageEvent) {
	for m := range c {
		switch m.Message.(type) {
		// TODO:
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

	// clientID := client.GetID()

	// TODO: Add the player to the game
}

func (mb *MonsterBash) removePlayer(client ws.HubClient) {
	// TODO: Remove the player from the game

}

func (mb *MonsterBash) initPlayer(clientID int32, name string) {
	// TODO: Initialize the player with name sent from the client
}
