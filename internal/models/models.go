package models

// Message es la estructura de cada mensaje del chat
type Message struct {
	Type     string `json:"type"`
	Content  string `json:"content"`
	Username string `json:"username"`
	Room     string `json:"room"`
}
