package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}

	diffs, errs := poll()

	commands := make(chan string)

	go func() {
		for {
			data := make(map[string]string)

			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					log.Printf("error: %v", err)
				}
				break
			}

			if err := json.Unmarshal(message, &data); err != nil {
				log.Printf("error: %v", err)
			}

			if cmd, ok := data["command"]; ok {
				commands <- cmd
			}
		}
	}()

	for {
		select {
		case diff := <-diffs:
			conn.WriteJSON(diff)
		case err := <-errs:
			conn.WriteJSON(gin.H{"error": err.Error()})
		case cmd := <-commands:
			client.Cmd(cmd)
		}
	}
}
