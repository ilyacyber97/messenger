package services

import (
	"context"
	"server/domain"

	"github.com/google/uuid"
)

type MessengerService struct {
	client domain.MessageServiceClient
}

func NewMessengerServiceRepository(client domain.MessageServiceClient) *MessengerService {
	return &MessengerService{client: client}
}

func (s *MessengerService) SaveMessage(message *domain.Message) error {
	message.Id = uuid.New().String()
	_, err := s.client.SaveMessage(context.Background(), message)
	return err
}

func (s *MessengerService) ReadMessage(id string) (*domain.Message, error) {
	request := &domain.ReadMessageRequest{Id: id}
	return s.client.ReadMessage(context.Background(), request)
}

func (s *MessengerService) ReadMessages() ([]*domain.Message, error) {
	empty := &domain.Empty{}
	list, err := s.client.ReadMessages(context.Background(), empty)
	return list.Messages, err

}
