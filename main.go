package main

import (
	"flag"
	"fmt"
	"monsterbash-server/pkg/config"
	"monsterbash-server/pkg/server"
)

func main() {
	defaultListenAddr := config.GetEnvVarOrDefault("MB_LISTEN_ADDR", ":3000")
	var addr = flag.String("addr", defaultListenAddr, "http server address")
	flag.Parse()

	s := server.NewServer(*addr)

	fmt.Println("Server listening on", *addr)

	if err := s.Start(); err != nil {
		panic(err)
	}

	fmt.Println("Server stopped")
}
