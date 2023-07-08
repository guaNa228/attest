package logger

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func Logger(logChan chan string, wsConn *websocket.Conn, writeToWs bool) {
	for message := range logChan {
		fmt.Println(message)
		if writeToWs {
			err := wsConn.WriteJSON(Message{Type: "log", Data: message})
			if err != nil {
				log.Printf("WebSocket write error (log): %v", err)
				return
			}
		}
	}
}

func ErrLogger(logChan chan error, errorCounter *int, wsConn *websocket.Conn, writeToWs bool) {
	for err := range logChan {
		*errorCounter++
		fmt.Println("error while parsing: ", err.Error())
		if writeToWs {
			err := wsConn.WriteJSON(Message{Type: "err", Data: err.Error()})
			if err != nil {
				log.Printf("WebSocket write error (log): %v", err)
				return
			}
		}
	}
}

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}
