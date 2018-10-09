package server

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/JCGrant/paint/statik"
	"github.com/gorilla/websocket"
	"github.com/rakyll/statik/fs"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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

func (s *Server) Serve() error {
	statikFS, err := fs.New()
	if err != nil {
		return fmt.Errorf("setting up statikFS failed: %s", err)
	}
	http.HandleFunc("/", http.FileServer(statikFS).ServeHTTP)
	http.HandleFunc("/ws", s.handleNewConnection)
	go s.hub.run()
	log.Printf("serving at http://127.0.0.1%s\n", s.address)
	err = http.ListenAndServe(s.address, nil)
	if err != nil {
		log.Fatal("serving failed: ", err)
	}
	return nil
}
