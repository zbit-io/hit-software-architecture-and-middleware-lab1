package broker

import (
	"github.com/zbit-io/hit-software-architecture-and-middleware-lab1/internal/message"
	"sync"
)

type Broker struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

// NewBroker 创建新的 Broker
func NewBroker() *Broker {
	return &Broker{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run 启动 Broker 的监听
func (b *Broker) Run() {
	for {
		select {
		case client := <-b.register:
			b.mu.Lock()
			b.clients[client] = true
			b.mu.Unlock()
		case client := <-b.unregister:
			b.mu.Lock()
			delete(b.clients, client)
			b.mu.Unlock()
		}
	}
}

func (b *Broker) Broadcast(msg message.Message) {
	b.mu.Lock()
	clients := make([]*Client, 0, len(b.clients))
	for client := range b.clients {
		clients = append(clients, client)
	}
	b.mu.Unlock()

	for _, client := range clients {
		select {
		case client.send <- msg:
			// 正常发送消息
		default:
			// 如果客户端的 `send` 通道阻塞，关闭该通道并从 `clients` 中删除客户端
			close(client.send)
			b.mu.Lock()
			delete(b.clients, client)
			b.mu.Unlock()
		}
	}
}
