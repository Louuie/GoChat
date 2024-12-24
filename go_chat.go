package main

import (
	"fmt"
	"gochat/fiber"
	"gochat/page"
	"gochat/ws"
)

func main() {
	// Start WebSocket server
	ws.StartServer(messageHandler)

	// Start fiber server
	fiber.StartServer()

	// Start the static HTML Page
	page.StartPage()

	// Block main thread to keep server running
	select {}
}

// Message handler processes messages from WebSocket clients
func messageHandler(message []byte) {
	fmt.Printf("Reccieved %s\n", message)
}
