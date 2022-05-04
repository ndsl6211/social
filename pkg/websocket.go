package pkg

import (
	"net/http"

	"github.com/gorilla/websocket"
)

func NewWebSocketUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
}
