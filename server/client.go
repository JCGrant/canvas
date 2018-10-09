package server

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeTimeout   = 10 * time.Second
	pongTimeout    = 60 * time.Second
	pingPeriod     = pongTimeout * 9 / 10
	maxNumMessages = 256
	maxMessageSize = 512
)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	messages chan []byte
}

func newClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:      hub,
		conn:     conn,
		messages: make(chan []byte, maxNumMessages),
	}
}

func (c *Client) send(message []byte) {
	select {
	case c.messages <- message:
	default:
		log.Println("closing client because it is dead")
		c.close()
		c.hub.unregister(c)
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister(c)
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	err := c.conn.SetReadDeadline(time.Now().Add(pongTimeout))
	if err != nil {
		log.Printf("setting read deadline failed: %s\n", err)
		return
	}
	c.conn.SetPongHandler(func(string) error {
		err := c.conn.SetReadDeadline(time.Now().Add(pongTimeout))
		if err != nil {
			return fmt.Errorf("setting read deadline failed: %s", err)
		}
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("reading message failed: %s\n", err)
			}
			break
		}
		c.hub.broadcast(message)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.messages:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err != nil {
				log.Printf("setting write deadline failed: %s\n", err)
				return
			}
			if !ok {
				log.Println("closing connection")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("getting next text message writer failed: %s\n", err)
				return
			}
			_, err = w.Write(message)
			if err != nil {
				log.Printf("writing message failed: %s\n", err)
				return
			}
			err = w.Close()
			if err != nil {
				log.Printf("closing writer failed: %s\n", err)
				return
			}
		case <-ticker.C:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err != nil {
				log.Printf("setting write deadline failed: %s\n", err)
				return
			}
			err = c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Printf("writing ping message failed: %s\n", err)
				return
			}
		}
	}
}

func (c *Client) close() {
	close(c.messages)
}
