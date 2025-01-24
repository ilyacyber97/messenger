package ports

import (
	"server/domain"
)

type Messenger interface {
	SaveMessage(message domain.Message) error
	ReadMessage(id string) (*domain.Message, error)
	ReadMessages() ([]*domain.Message, error)
}
