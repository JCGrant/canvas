package server

import "log"

type Hub struct {
	clients      map[*Client]bool
	registered   chan *Client
	unregistered chan *Client
	messages     chan []byte
}

func newHub() *Hub {
	return &Hub{
		clients:      make(map[*Client]bool),
		registered:   make(chan *Client),
		unregistered: make(chan *Client),
		messages:     make(chan []byte),
	}
}

func (h *Hub) register(client *Client) {
	h.registered <- client
}

func (h *Hub) unregister(client *Client) {
	h.unregistered <- client
}

func (h *Hub) broadcast(message []byte) {
	log.Printf("broadcasting: '%s'", message)
	h.messages <- message
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.registered:
			h.clients[client] = true
		case client := <-h.unregistered:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.close()
			}
		case message := <-h.messages:
			for client := range h.clients {
				client.send(message)
			}
		}
	}
}
