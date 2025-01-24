package main

import (
	"fmt"
	"os"
	"server/internal/adapters/domain"
	"server/internal/adapters/http"
	"server/log"
	"server/metrics/prometric"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var handler *http.HTTPHandler

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
	//log
	logger, err := initLog()
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	//metric
	prometheus.MustRegister(prometric.RequestCounter)

	clientConn, err := grpc.NewClient("server1:50054", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("Не удалось подключиться: %v", zap.Error(err))
	}
	defer clientConn.Close()

	client := domain.NewMessageServiceClient(clientConn)

	handler = http.NewHandlerMessengerServiceRepository(svc)

	router := gin.Default()
	router.POST("message", handler.SaveMessage)
	router.GET("message/:id", handler.ReadMessage)
	router.GET("/messages", handler.ReadMessages)
	//metrics
	router.GET("/metrics", prometric.MetricHandler)
	logger.Info("Запуск сервера на порту :8080")
	router.Run(":8080")

}
