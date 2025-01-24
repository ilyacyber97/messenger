package services

import (
	"server/internal/core/domain"
	"server/internal/core/ports"

	"github.com/google/uuid"
)

type messengerService struct {
	repo ports.MessengerRepository
}

func NewMessengerServiceRepository(repo ports.MessengerRepository) ports.MessengerService {
	return &messengerService{repo: repo}
}

func (s *messengerService) SaveMessage(message domain.Message) error {
	message.Id = uuid.New().String()
	return s.repo.SaveMessage(message)
}

func (s *messengerService) ReadMessage(id string) (*domain.Message, error) {
	return s.repo.ReadMessage(id)
}
func (s *messengerService) ReadMessages() ([]*domain.Message, error) {
	return s.repo.ReadMessages()
}
