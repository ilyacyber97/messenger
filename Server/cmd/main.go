package main

import (
	"flag"
	"messenger/config"
	"messenger/internal/adapters/handler/http"
	"messenger/internal/adapters/repository/postgres"
	"messenger/internal/adapters/repository/redis"
	"messenger/internal/core/services"
	"messenger/logs/zaplog"
	"messenger/metrics/prometric"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gin-gonic/gin"
)

var flagDB = flag.String("db", "postgres", "Выберита Базу Данных")
var handler *http.HTTPHandler

func init() {
	prometheus.MustRegister(prometric.RequestCounter)
	zaplog.Init()
}

func main() {
	//log
	defer zaplog.Shutdown()

	//flag
	flag.Parse()

	switch *flagDB {

	case "redis":
		zaplog.Logger.Info("Пользователь выбрал Redis")

		repo, err := redis.NewRedisRepository(config.RedisHost)
		if err != nil {
			zaplog.Logger.Fatal("Не удалось подключиться к Redis", zaplog.LogError(err))
		}

		svc := services.NewMessengerServiceRepository(repo)
		handler = http.NewHandlerMessengerServiceRepository(svc)

	default:
		zaplog.Logger.Info("Пользователь выбрал Postgresql")

		repo, err := postgres.NewPostgresRepository(config.PostgresDSN)
		if err != nil {
			zaplog.Logger.Fatal("Не удалось подключиться к PostgreSQL", zaplog.LogError(err))
		}

		svc := services.NewMessengerServiceRepository(repo)
		handler = http.NewHandlerMessengerServiceRepository(svc)

	}

	router := gin.Default()

	router.POST("message", handler.SaveMessage)
	router.GET("message/:id", handler.ReadMessage)
	router.GET("/messages", handler.ReadMessages)
	//metrics
	router.GET("/metrics", prometric.MetricHandler)
	zaplog.Logger.Info("Запуск сервера на порту :8080")
	router.Run(":8080")

}
