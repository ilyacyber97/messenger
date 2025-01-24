package main

import (
	"flag"
	"server/config"
	"server/internal/adapters/handler/http"
	"server/internal/adapters/repository/postgres"
	"server/internal/adapters/repository/redis"
	"server/internal/core/services"
	"server/logs/zaplog"
	"server/metrics/prometric"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gin-gonic/gin"
)

var handler *http.HTTPHandler

func init() {
	prometheus.MustRegister(prometric.RequestCounter)
	zaplog.Init()
}

func main() {
			
		repo, err := postgres.NewPostgresRepository(config.PostgresDSN)
		if err != nil {
			zaplog.Logger.Fatal("Не удалось подключиться к PostgreSQL", zaplog.LogError(err))
		}

		svc := services.NewMessengerServiceRepository(repo)
		handler = http.NewHandlerMessengerServiceRepository(svc)

	

	router := gin.Default()

	router.POST("message", handler.SaveMessage)
	router.GET("message/:id", handler.ReadMessage)
	router.GET("/messages", handler.ReadMessages)
	//metrics
	router.GET("/metrics", prometric.MetricHandler)
	zaplog.Logger.Info("Запуск сервера на порту :8080")
	router.Run(":8080")

}
