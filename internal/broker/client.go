package broker

import (
	"github.com/gorilla/websocket"
	"github.com/zbit-io/hit-software-architecture-and-middleware-lab1/internal/message"
)

type Client struct {
	conn *websocket.Conn
	send chan message.Message
}
