package hub

import "chat-app/internal/models"

// Hub es el centro que gestiona todos los clientes conectados
type Hub struct {
	// Mapa de salas con sus clientes
	Rooms map[string]map[*Client]bool
	// Canal para registrar nuevos clientes
	Register chan *Client
	// Canal para desconectar clientes
	Unregister chan *Client
	// Canal para distribuir mensajes
	Broadcast chan models.Message
}

// NewHub crea un nuevo Hub
func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan models.Message),
	}
}

// Run es el corazón del Hub — corre en su propia goroutine
func (h *Hub) Run() {
	for {
		select {

		// Nuevo cliente se conecta
		case client := <-h.Register:
			if _, ok := h.Rooms[client.Room]; !ok {
				h.Rooms[client.Room] = make(map[*Client]bool)
			}
			h.Rooms[client.Room][client] = true

		// Cliente se desconecta
		case client := <-h.Unregister:
			if _, ok := h.Rooms[client.Room]; ok {
				delete(h.Rooms[client.Room], client)
				close(client.Send)
			}

		// Distribuir mensaje a todos en la sala
		case message := <-h.Broadcast:
			if clients, ok := h.Rooms[message.Room]; ok {
				for client := range clients {
					client.Send <- message
				}
			}
		}
	}
}
