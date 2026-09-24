package server

import (
	"fmt"
	"net/http"

	"discord/internal/websocket"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// type Server struct {
// 	db         *sql.DB
// 	repository *Repository
// 	service    *Service
// 	handler    *Handler
// }

func NewServer(db *pgxpool.Pool) *http.Server {
	repository := NewRepository(db)
	hub := websocket.NewHub()
	service := NewService(repository, hub)
	handler := NewHandler(service)
	wsHandler := websocket.NewHandler(hub)

	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {

		r.Get("/users", handler.GetUsers)
		r.Post("/login", handler.Authorization)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, World!")
		})

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			r.Get("/auth/me", handler.Me)
			r.Get("/servers", handler.GetServers)
			r.Get("/server/{server_id}/channel/{channel_id}/chat/{chat_id}/messages", handler.GetMessages)
			r.Post("/server/{server_id}/channel/{channel_id}/chat/{chat_id}/message", handler.CreateMessage)
		})
	})
	r.Route("/", func(r chi.Router) {
		// r.Get("/ws", handler.wsHandler)
		r.Get("/ws", wsHandler.Handle)
	})

	go hub.Run()

	return &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: r,
	}
}
