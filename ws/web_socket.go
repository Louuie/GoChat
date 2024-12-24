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
	clients       map[*websocket.Conn]bool
	UUID          map[string]string
	handleMessage func(message []byte)
}

func StartServer(handleMessage func(message []byte)) *Server {
	server := Server{
		make(map[*websocket.Conn]bool),
		make(map[string]string),
		handleMessage,
	}
	http.HandleFunc("/ws", server.echo)
	go http.ListenAndServe(":8080", nil)
	return &server
}

func (server *Server) echo(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln(err)
	}
	var chatUUID = "42fsqg-kl42fe"
	server.UUID[chatUUID] = "1"
	server.clients[conn] = true
	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			log.Fatalln(err)
			break
		}
		go server.handleMessage(msg)
		go server.WriteMessage(msg)
	}
	conn.Close()
	delete(server.clients, conn)
	delete(server.UUID, chatUUID)
}

func (server *Server) WriteMessage(message []byte) {
	for conn := range server.clients {
		conn.WriteMessage(websocket.TextMessage, message)
	}
}
