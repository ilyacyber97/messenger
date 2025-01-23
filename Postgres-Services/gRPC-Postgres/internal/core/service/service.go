package service

import (
	"grpc-postgres/domain"
	"grpc-postgres/internal/core/ports"
)

type service struct {
	repo ports.MessageService
}

func NewServiceRepository(repo ports.MessageService) *service {
	return &service{repo: repo}
}

func (s *service) SaveMessage(message *domain.Message) error {
	return s.repo.SaveMessage(message)
}

func (s *service) ReadMessage(id string) (*domain.Message, error) {
	return s.repo.ReadMessage(id)
}
func (s *service) ReadMessages() ([]*domain.Message, error) {
	return s.repo.ReadMessages()
}
