package main

import (
	"fmt"
	"grpc-postgres/config"
	"grpc-postgres/domain"
	"grpc-postgres/internal/adapters/repository/postgres"
	"grpc-postgres/internal/adapters/server"
	"grpc-postgres/internal/core/service"
	"grpc-postgres/log"
	"net"
	"os"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func initLog() (*log.Logger, error) {
	file, err := os.Create("../logfile.txt")
	if err != nil {
		return nil, fmt.Errorf("модуль zap: ошибка создания log файла %v", err)
	}
	logger := log.NewLoggerZap(file)

	logger.Info("Logger инициализирован успешно")
	return logger, nil
}

func main() {
	// логирование
	logger, err := initLog()
	defer logger.Close()

	repo, err := postgres.NewPostgresRepository(config.DSN, logger)
	service := service.NewServiceRepository(repo)
	server := server.NewServer(logger, service)

	// Создаем gRPC сервер
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}
	s := grpc.NewServer()

	domain.RegisterMessageServiceServer(s, server)

	// Запускаем сервер
	if err := s.Serve(lis); err != nil {
		logger.Fatal("failed to serve: ", zap.Error(err))
	}

}
