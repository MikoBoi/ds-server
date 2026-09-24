package server

import (
	"context"
	"encoding/json"

	"discord/internal/model"
	"discord/internal/websocket"
)

type Service struct {
	repository *Repository
	hub        *websocket.Hub
}

func NewService(repository *Repository, hub *websocket.Hub) *Service {
	return &Service{repository: repository, hub: hub}
}

func (s *Service) GetServers(ctx context.Context, userID int64) ([]model.Server, error) {
	return s.repository.GetServers(ctx, userID)
}

func (s *Service) GetMessages(ctx context.Context, server_id, channel_id, chat_id int64) ([]model.Message, error) {
	return s.repository.GetMessages(ctx, server_id, channel_id, chat_id)
}

func (s *Service) GetUsers(ctx context.Context) ([]model.User, error) {
	return s.repository.GetUsers(ctx)
}

func (s *Service) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	return s.repository.GetUser(ctx, userID)
}

func (s *Service) Authorization(ctx context.Context, email, password string) (string, error) {
	return s.repository.Authorization(ctx, email, password)
}

func (s *Service) CreateMessage(ctx context.Context, server_id, channel_id, chat_id, userID int64, content string) error {
	err := s.repository.CreateMessage(ctx, server_id, channel_id, chat_id, userID, content)
	if err != nil {
		return err
	}

	event := model.MessageEvent{
		Type: "message.created",
		Message: model.Message{
			ServerId:  int(server_id),
			ChannelId: int(channel_id),
			ChatId:    int(chat_id),
			AuthorId:  int(userID),
			Content:   content,
		},
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	s.hub.Broadcast(data)

	return nil
}
