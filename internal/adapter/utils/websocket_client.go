package utils

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	UserId uuid.UUID
	Conn   *websocket.Conn
}

func NewWebSocketClient(userId uuid.UUID, conn *websocket.Conn) *WebSocketClient {
	return &WebSocketClient{userId, conn}
}
