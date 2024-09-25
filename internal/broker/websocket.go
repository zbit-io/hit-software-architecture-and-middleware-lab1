package broker

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zbit-io/hit-software-architecture-and-middleware-lab1/internal/message"
	"net/http"
)

// WebSocket 处理连接
func HandleWebSocket(broker *Broker, w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	client := &Client{conn: conn, send: make(chan message.Message)}
	broker.register <- client

	go func() {
		defer func() {
			broker.unregister <- client
			conn.Close()
		}()

		for {
			var msg message.Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				fmt.Println("Error reading message:", err)
				break
			}
			broker.Broadcast(msg)
		}
	}()

	for msg := range client.send {
		err := conn.WriteJSON(msg)
		if err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}
}
