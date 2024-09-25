package main

import (
	"fmt"
	"github.com/zbit-io/hit-software-architecture-and-middleware-lab1/internal/broker"
	"net/http"
)

func main() {
	currBroker := broker.NewBroker()
	go currBroker.Run()
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		broker.HandleWebSocket(currBroker, w, r)
	})

	fmt.Println("WebSocket server started at ws://localhost:8080/chat")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
