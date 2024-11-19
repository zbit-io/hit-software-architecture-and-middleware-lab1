package broker

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zbit-io/hit-software-architecture-and-middleware-lab1/internal/message"
	"net/http"
	"time"
)

// WebSocket 处理连接
func HandleWebSocket(broker *Broker, w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	client := &Client{conn: conn, send: make(chan message.Message, 100)}
	broker.register <- client

	// 启动单独的 goroutine 处理发送消息，保证写操作是串行的
	go client.writePump()

	// 设置读超时和 pong handler
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second)) // 每次收到 pong 时重置超时
		return nil
	})

	// 启动心跳 Ping 机制
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					fmt.Println("Ping error:", err)
					broker.HandleWebSocketClose(conn, client)
					return
				}
			}
		}
	}()

	// 处理接收消息的 goroutine
	go func() {
		defer func() {
			broker.HandleWebSocketClose(conn, client)
		}()

		for {
			var msg message.Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Printf("WebSocket closed unexpectedly: %v\n", err)
				} else {
					fmt.Println("Error reading message:", err)
				}
				break
			}
			broker.Broadcast(msg)
		}
	}()
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for msg := range c.send {
		if err := c.conn.WriteJSON(msg); err != nil {
			fmt.Println("Error writing message:", err)
			return
		}
	}
}

func (b *Broker) HandleWebSocketClose(conn *websocket.Conn, client *Client) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.clients[client]; !ok {
		return // 已经从客户端列表中移除，不再处理
	}
	// 从 Broker 中移除客户端
	delete(b.clients, client)

	// 关闭客户端的 send 通道
	close(client.send)

	// 发送关闭消息
	if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		fmt.Printf("Error sending close message: %v\n", err)
	}
	conn.Close()
}
