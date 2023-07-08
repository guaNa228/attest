package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var GlobalWsConn *websocket.Conn
var GlobalWsWg sync.WaitGroup

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Customize this according to your security requirements
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 3000,
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	GlobalWsWg = sync.WaitGroup{}

	GlobalWsWg.Add(1)

	GlobalWsConn = conn

	GlobalWsWg.Wait()

	time.Sleep(10 * time.Second)

	GlobalWsConn.Close()
}
