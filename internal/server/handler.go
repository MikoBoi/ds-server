package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetServers(w http.ResponseWriter, r *http.Request) {
	servers, err := h.service.GetServers(r.Context(), r.Context().Value("user_id").(int64))
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(servers)
}

func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	serverID, err := strconv.ParseInt(chi.URLParam(r, "server_id"), 10, 64)
	channelID, err := strconv.ParseInt(chi.URLParam(r, "channel_id"), 10, 64)
	chatID, err := strconv.ParseInt(chi.URLParam(r, "chat_id"), 10, 64)

	messages, err := h.service.GetMessages(r.Context(), serverID, channelID, chatID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) Authorization(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	loginResponse, err := h.service.Authorization(r.Context(), loginRequest.Email, loginRequest.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse)
}

func (h *Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Message string `json:"message"`
	}

	serverID, err := strconv.ParseInt(chi.URLParam(r, "server_id"), 10, 64)
	channelID, err := strconv.ParseInt(chi.URLParam(r, "channel_id"), 10, 64)
	chatID, err := strconv.ParseInt(chi.URLParam(r, "chat_id"), 10, 64)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = h.service.CreateMessage(r.Context(), serverID, channelID, chatID, r.Context().Value("user_id").(int64), req.Message)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
