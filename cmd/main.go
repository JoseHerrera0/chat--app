package main

import (
	"chat-app/internal/handlers"
	"chat-app/internal/hub"
	"fmt"
	"net/http"
)

func main() {
	h := hub.NewHub()
	go h.Run()

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.ServeWS(h, w, r)
	})

	fmt.Println("Servidor corriendo en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
