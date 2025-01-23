package ports

import (
	"grpc-postgres/domain"

	"go.uber.org/zap"
)

type MessageService interface {
	SaveMessage(*domain.Message) error
	ReadMessage(string) (*domain.Message, error)
	ReadMessages() ([]*domain.Message, error)
}

type LoggerPort interface {
	Close()
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Fd() uintptr
	Info(msg string, fields ...zap.Field)
}
