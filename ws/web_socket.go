package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	Rooms          map[string]map[*websocket.Conn]bool
	upgrader       websocket.Upgrader
	handleMessages func(message []byte)
}

func StartServer(handleMessage func(message []byte)) *Server {
	server := &Server{
		Rooms:          make(map[string]map[*websocket.Conn]bool), //  map that returns whether the room number is attachted or a part of the connection
		upgrader:       upgrader,
		handleMessages: handleMessage,
	}
	http.HandleFunc("/ws", server.echo)
	go http.ListenAndServe(":8080", nil)
	return server
}

func (server *Server) echo(w http.ResponseWriter, r *http.Request) {
	conn, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	// Get parameters from the URL
	params := r.URL.Query()
	chatRoomNumber := params.Get("room")
	chatUUID := params.Get("uuid")

	if chatRoomNumber == "" || chatUUID == "" {
		conn.WriteMessage(websocket.CloseMessage, []byte("Invalid room or UUID"))
		conn.Close()
		return
	}

	// Add client to the room
	if _, exists := server.Rooms[chatRoomNumber]; !exists {
		server.Rooms[chatRoomNumber] = make(map[*websocket.Conn]bool)
	}
	server.Rooms[chatRoomNumber][conn] = true
	log.Printf("Client with UUID %s joined room %s", chatUUID, chatRoomNumber)
	log.Printf("Current state of Rooms: %+v\n", server.Rooms)

	// Remove client from room on disconnect
	defer func() {
		conn.Close()
		delete(server.Rooms[chatRoomNumber], conn)
		if len(server.Rooms[chatRoomNumber]) == 0 {
			delete(server.Rooms, chatRoomNumber) // Clean up empty room
		}
		log.Printf("Client with UUID %s disconnected from room %s", chatUUID, chatRoomNumber)
		log.Printf("Current state of Rooms: %+v\n", server.Rooms)
	}()

	// Read messages and broadcast to room
	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client: %v\n", err)
			break
		}

		log.Printf("Message received in room %s: %s", chatRoomNumber, string(msg))

		// Broadcast the message to all clients in the same room
		for client := range server.Rooms[chatRoomNumber] {
			if client != conn { // Exclude the sender if desired
				err := client.WriteMessage(mt, msg)
				if err != nil {
					log.Printf("Error sending message to client: %v\n", err)
					client.Close()
					delete(server.Rooms[chatRoomNumber], client)
				}
			}
		}
	}
}
