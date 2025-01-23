package redis

import (
	"encoding/json"
	"fmt"
	"messenger/config"
	"messenger/internal/core/domain"
	"messenger/logs/zaplog"
	"strconv"

	"github.com/go-redis/redis/v7"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(host string) (*RedisRepository, error) {
	zaplog.Logger.Info("Подключение к  Redis", zaplog.LogString("dsn", config.PostgresDSN))

	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "",
		DB:       0,
	})
	// Проверка подключения
	_, err := client.Ping().Result()
	if err != nil {

		zaplog.Logger.Error("модуль redis: Ошбка подключения к Redis", zaplog.LogError(err))

		return nil, fmt.Errorf("модуль redis: ошибка подключения к Redis: %w", err)
	}
	return &RedisRepository{client: client}, nil
}

func (r *RedisRepository) SaveMessage(message domain.Message) error {
	zaplog.Logger.Info("модуль redis: Сохранение message", zaplog.LogString("id", message.Id),
		zaplog.LogString("body", message.Body))
	json, err := json.Marshal(message)
	if err != nil {
		zaplog.Logger.Error("модуль redis: Ошибка сохранения message", zaplog.LogString("id", message.Id), zaplog.LogError(err))
		return fmt.Errorf("модуль redis: не удалось сохранить сообщение с id %s: %w", message.Id, err)
	}
	r.client.HSet("messages", message.Id, json)
	zaplog.Logger.Info("модуль redis: Успешная Запись message", zaplog.LogString("id", message.Id), zaplog.LogString("body", message.Body))

	return nil
}

func (r *RedisRepository) ReadMessage(id string) (*domain.Message, error) {
	zaplog.Logger.Info("модуль redis: Чтение message", zaplog.LogString("id", id))

	value, err := r.client.HGet("messages", id).Result()
	if err != nil {
		zaplog.Logger.Error("модуль redis: Ошибка чтения message", zaplog.LogString("id", id), zaplog.LogError(err))

		return nil, fmt.Errorf("модуль redis: не удалось прочитать сообщение с id %s: %w", id, err)
	}

	message := &domain.Message{}
	err = json.Unmarshal([]byte(value), message)
	if err != nil {
		return nil, err
	}
	zaplog.Logger.Info("модуль redis: Успешное Чтение message", zaplog.LogString("id", message.Id), zaplog.LogString("body", message.Body))

	return message, nil
}

func (r *RedisRepository) ReadMessages() ([]*domain.Message, error) {
	zaplog.Logger.Info("модуль redis: Чтение messages")

	messages := []*domain.Message{}
	maps, err := r.client.HGetAll("messages").Result()
	if err != nil {
		zaplog.Logger.Error("модуль redis: Ошибка чтения messages")

		return nil, fmt.Errorf("модуль redis: не удалось прочитать сообщения: %w", err)

	}
	fmt.Println(maps)
	for _, v := range maps {
		message := &domain.Message{}
		err := json.Unmarshal([]byte(v), message)
		if err != nil {
			return nil, fmt.Errorf("модуль redis: не удалось прочитать сообщения: %w", err)

		}
		messages = append(messages, message)
	}
	zaplog.Logger.Info("модуль redis: Messages прочитаны успешно", zaplog.LogString("number", strconv.Itoa(len(maps))))

	return messages, nil

}
