package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Sirve los archivos de la carpeta static
	http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Println("Servidor corriendo en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
