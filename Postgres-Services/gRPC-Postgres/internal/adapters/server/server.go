package server

import (
	"context"
	"fmt"
	"grpc-postgres/domain"
	"grpc-postgres/internal/core/ports"
	"grpc-postgres/log"
	"strconv"

	"go.uber.org/zap"
)

type server struct {
	log *log.Logger
	domain.UnimplementedMessageServiceServer
	service ports.MessageService
}

func NewServer(log *log.Logger, service ports.MessageService) *server {

	return &server{log: log, service: service}
}

func (s *server) SaveMessage(ctx context.Context, message *domain.Message) (*domain.Empty, error) {
	s.log.Info("server: Получен запрос SaveMessage", zap.String(message.Id, message.Body))
	err := s.service.SaveMessage(message)
	if err != nil {
		s.log.Error("server: Ошибка сохранения message", zap.Error(err))
		return nil, fmt.Errorf("server: Ошибка сохранения message: %w", err)
	}
	s.log.Info("server: Message сохранен успешно", zap.String("message", message.Body))
	return &domain.Empty{}, nil
}

func (s *server) ReadMessage(ctx context.Context, req *domain.ReadMessageRequest) (*domain.Message, error) {
	s.log.Info("server: Получен запрос ReadMessage", zap.String("client_ip", req.Id))

	message, err := s.service.ReadMessage(req.Id)
	if err != nil {
		s.log.Error("server: Ошибка чтения message", zap.String("id", req.Id), zap.Error(err))
		return nil, fmt.Errorf("server: Ошибка чтения message по id %w: %w", req.Id, err)
	}

	s.log.Info("server: Message прочитан успешно", zap.String("id", req.Id))

	return message, nil

}

func (s *server) ReadMessages(ctx context.Context, req *domain.Empty) (*domain.MessagesList, error) {
	s.log.Info("server: Получен запрос ReadMessages")
	slice, err := s.service.ReadMessages()
	if err != nil {
		s.log.Error("server: Ошибка чтения messages", zap.Error(err))
		return nil, fmt.Errorf("server: Ошибка чтения messages: %w", err)

	}
	s.log.Info("server: Messages прочитан успешно", zap.String("number", strconv.Itoa(len(slice))))
	return &domain.MessagesList{Messages: slice}, nil
}
