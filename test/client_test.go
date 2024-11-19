// broker_test.go
package test

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zbit-io/hit-software-architecture-and-middleware-lab1/internal/message"
	"net/url"
	"sync"
	"testing"
)

func BenchmarkMultipleWebSocketClients(b *testing.B) {

	// WebSocket 服务器的URL
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/chat"}

	var mu sync.Mutex
	var errors []error

	// 重置计时器，忽略启动服务器的时间
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 建立 WebSocket 连接
			conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("WebSocket dial error: %v", err))
				mu.Unlock()
				continue // 记录错误后继续下一个迭代
			}

			// 发送和接收消息
			for i := 0; i < 10; i++ {
				msg := message.NewMessage("user1", "user2", "hello")
				if err := conn.WriteJSON(msg); err != nil {
					mu.Lock()
					errors = append(errors, fmt.Errorf("Write message error: %v", err))
					mu.Unlock()
					break // 发生错误，退出消息循环
				}

				_, _, err = conn.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						mu.Lock()
						errors = append(errors, fmt.Errorf("WebSocket connection closed unexpectedly: %v", err))
						mu.Unlock()
					} else {
						mu.Lock()
						errors = append(errors, fmt.Errorf("Read message error: %v", err))
						mu.Unlock()
					}
					break // 发生错误，退出消息循环
				}
			}

			if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil && !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				mu.Lock()
				errors = append(errors, fmt.Errorf("Error sending close message: %v", err))
				mu.Unlock()
			}
			conn.Close()

		}
	})

	b.StopTimer()

	// 处理记录的错误
	if len(errors) > 0 {
		for _, err := range errors {
			b.Error(err)
		}
	}
}
