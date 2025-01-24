package services

import (
	"server/domain"
	"server/internal/core/ports"
)

type messengerService struct {
	log    ports.MessengerLog
	client domain.MessageServiceClient
}

func NewMessengerServiceRepository(log ports.MessengerLog, client domain.MessageServiceClient) *messengerService {
	return &messengerService{log: log, client: client}
}

func (s *messengerService) SaveMessage(message domain.Message) error {
	s.client.SaveMessage()
	return s.repo.SaveMessage(message)
}

func (s *messengerService) ReadMessage(id string) (*domain.Message, error) {
	return s.repo.ReadMessage(id)
}
func (s *messengerService) ReadMessages() ([]*domain.Message, error) {
	return s.repo.ReadMessages()
}
