package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"monsterbash-server/pkg/monsterbash"
	"monsterbash-server/pkg/ws"
	"net/http"
)

var BufferSize = 1024

type Server struct {
	listenAddr string
	hub        *ws.Hub
	game       *monsterbash.MonsterBash
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		hub:        ws.NewHub(),
		game:       monsterbash.NewMonsterBash(),
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/health", s.handleHealthRoute)
	http.HandleFunc("/ws", s.handleWsRoute)

	go s.hub.Run()
	go s.game.Run()

	// Game processing messages from the websocket hub
	go s.game.ProcessNewPlayers(s.hub.GetPlayerJoinGameChannel())
	go s.game.ProcessUnregisteredPlayers(s.hub.GetUnregisterChannel())
	go s.game.ProcessInput(s.hub.GetMessageReceiveChannel())

	// Websocket hub processing messages from the game
	go s.hub.ProcessSendMessages(s.game.GetSendMessageChannel())
	go s.hub.ProcessBroadcastMessages(s.game.GetBroadcastMessageChannel())

	return http.ListenAndServe(s.listenAddr, nil)
}

func (s *Server) handleHealthRoute(w http.ResponseWriter, _ *http.Request) {
	responseData := map[string]bool{"ok": true}

	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		log.Println("Error while processing /health request:", err.Error())
	}
}

func (s *Server) handleWsRoute(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  BufferSize,
		WriteBufferSize: BufferSize,
	}

	// Allow requests from approved origins
	upgrader.CheckOrigin = func(r *http.Request) bool {
		// TODO: Move this to a config file
		acceptedOrigins := []string{"http://localhost:3000"}
		for _, origin := range acceptedOrigins {
			if r.Header.Get("Origin") == origin {
				return true
			}
		}

		return false
	}

	// Upgrade initial GET request to a websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Register the new client with the hub
	s.hub.Register(conn)
}
