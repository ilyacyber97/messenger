package ports

import (
	"server/domain"
)

type MessengerClient interface {
	SaveMessage(message domain.Message) error
	ReadMessage(id string) (*domain.Message, error)
	ReadMessages() ([]*domain.Message, error)
}

type MessengerLog interface {
	Close()
}
