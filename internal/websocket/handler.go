package websocket

import (
	"log"
	"net/http"

	"github.com/coder/websocket"
)

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return
	}

	defer conn.CloseNow()

	client := NewClient(conn)

	h.hub.register <- client

	log.Println("WebSocket client connected")

	for {
		_, _, err := conn.Read(r.Context())
		if err != nil {
			break
		}
	}

	h.hub.unregister <- client

	log.Println("WebSocket client disconnected")
}
