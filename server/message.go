package server

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// WSMessage define a basic Websocket message that has type and data
type WSMessage struct {
	MessageType int
	Data        []byte
}

// String implements
func (m *WSMessage) String() string {
	return fmt.Sprintf("Type %d, Data: %s)", m.MessageType, string(m.Data))
}

// NewWSTextMessage helper function to simplify creation of text messages
func NewWSTextMessage(bs []byte) *WSMessage {
	return &WSMessage{MessageType: websocket.TextMessage, Data: bs}
}

// NewWSMessage creates new WSMessage
func NewWSMessage(messageType int, bs []byte) *WSMessage {
	return &WSMessage{
		Data:        bs,
		MessageType: messageType,
	}
}
