package ports

import "messenger/internal/core/domain"

type MessengerService interface {
	SaveMessage(message domain.Message) error
	ReadMessage(id string) (*domain.Message, error)
	ReadMessages() ([]*domain.Message, error)
}
type MessengerRepository interface {
	SaveMessage(message domain.Message) error
	ReadMessage(id string) (*domain.Message, error)
	ReadMessages() ([]*domain.Message, error)
}
