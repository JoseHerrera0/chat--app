package hub

import (
	"chat-app/internal/models"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = 54 * time.Second
)

// Client representa un usuario conectado
type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan models.Message
	Username string
	Room     string
}

// ReadPump escucha mensajes que llegan desde el navegador
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, rawMsg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg models.Message
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			log.Println("error leyendo mensaje:", err)
			continue
		}

		msg.Username = c.Username
		msg.Room = c.Room

		c.Hub.Broadcast <- msg
	}
}

// WritePump envía mensajes al navegador
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			data, err := json.Marshal(msg)
			if err != nil {
				return
			}
			c.Conn.WriteMessage(websocket.TextMessage, data)

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
