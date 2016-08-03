package main

import (
	"fmt"
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

	for {
		select {
		case diff := <-diffs:
			conn.WriteJSON(diff)
		case err := <-errs:
			conn.WriteJSON(gin.H{"error": err.Error()})
		}
	}
}
