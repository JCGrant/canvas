package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Server struct {
	address string
	hub     *Hub
}

func New(address string) *Server {
	hub := newHub()
	return &Server{
		address: address,
		hub:     hub,
	}
}

func (s *Server) handleNewConnection(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrading connection failed: %s\n", err)
		return
	}
	client := newClient(s.hub, conn)
	s.hub.register(client)
	go client.readPump()
	go client.writePump()
}

func (s *Server) Serve() {
	go s.hub.run()
	http.HandleFunc("/ws", s.handleNewConnection)
	log.Printf("serving at http://127.0.0.1%s\n", s.address)
	err := http.ListenAndServe(s.address, nil)
	if err != nil {
		log.Fatal("serving failed: ", err)
	}
}
