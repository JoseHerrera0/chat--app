package handlers

import (
	"chat-app/internal/hub"
	"chat-app/internal/models"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ServeWS maneja la conexión WebSocket
func ServeWS(h *hub.Hub, w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	room := r.URL.Query().Get("room")

	if username == "" || room == "" {
		http.Error(w, "username y room son requeridos", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "error al conectar WebSocket", http.StatusInternalServerError)
		return
	}

	client := &hub.Client{
		Hub:      h,
		Conn:     conn,
		Send:     make(chan models.Message, 256),
		Username: username,
		Room:     room,
	}

	h.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
