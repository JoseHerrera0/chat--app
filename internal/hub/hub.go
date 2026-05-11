package hub

import (
	"chat-app/internal/models"
	"fmt"
)

type Hub struct {
	Rooms      map[string]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan models.Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan models.Message),
	}
}

// UsernameExiste verifica si un nombre ya está en uso en una sala
func (h *Hub) UsernameExiste(username, room string) bool {
	if clients, ok := h.Rooms[room]; ok {
		for client := range clients {
			if client.Username == username {
				return true
			}
		}
	}
	return false
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if _, ok := h.Rooms[client.Room]; !ok {
				h.Rooms[client.Room] = make(map[*Client]bool)
			}
			h.Rooms[client.Room][client] = true

			fmt.Println("Cliente registrado:", client.Username, "en sala:", client.Room)

			go func() {
				h.Broadcast <- models.Message{
					Type:     "join",
					Content:  client.Username + " se unió a la sala",
					Username: "Sistema",
					Room:     client.Room,
				}
			}()

		case client := <-h.Unregister:
			if _, ok := h.Rooms[client.Room]; ok {
				delete(h.Rooms[client.Room], client)
				close(client.Send)

				if len(h.Rooms[client.Room]) == 0 {
					delete(h.Rooms, client.Room)
				}

				go func() {
					h.Broadcast <- models.Message{
						Type:     "leave",
						Content:  client.Username + " salió de la sala",
						Username: "Sistema",
						Room:     client.Room,
					}
				}()
			}

		case message := <-h.Broadcast:
			if clients, ok := h.Rooms[message.Room]; ok {
				for client := range clients {
					select {
					case client.Send <- message:
					default:
						delete(h.Rooms[message.Room], client)
						close(client.Send)
					}
				}
			}
		}
	}
}
