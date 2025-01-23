package postgres

import (
	"fmt"
	"messenger/config"
	"messenger/internal/core/domain"
	"messenger/logs/zaplog"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type MessengerPostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(DSN string) (*MessengerPostgresRepository, error) {
	// Подключение к БД
	zaplog.Logger.Info("модуль postgres: Подключение к  PostgreSQL", zaplog.LogString("dsn", DSN))

	db, err := sqlx.Connect("postgres", config.PostgresDSN)
	if err != nil {
		zaplog.Logger.Error("модуль postgres: Ошбка подключения к PostgreSQL", zaplog.LogError(err))
		return nil, fmt.Errorf("модуль postgres: ошибка подключения к PostgreSQL: %w", err)
	}

	zaplog.Logger.Info("модуль postgres: Успешное подключение к PostgreSQL", zaplog.LogString("DSN", config.PostgresDSN))

	return &MessengerPostgresRepository{db: db}, nil
}

func (r *MessengerPostgresRepository) SaveMessage(message domain.Message) error {
	zaplog.Logger.Info("модуль postgres: Сохранение message", zaplog.LogString("id", message.Id),
		zaplog.LogString("body", message.Body))

	query := "INSERT INTO messenger(id ,body) VALUES(:id, :body)"
	_, err := r.db.NamedExec(query, message)
	if err != nil {
		zaplog.Logger.Error("модуль postgres: Ошибка сохранения message", zaplog.LogString("id", message.Id), zaplog.LogError(err))
		return fmt.Errorf("модуль postgres: не удалось сохранить сообщение с id %s: %w", message.Id, err)

	}

	zaplog.Logger.Info("модуль postgres: Успешная Запись message", zaplog.LogString("id", message.Id), zaplog.LogString("body", message.Body))

	return err
}

func (r *MessengerPostgresRepository) ReadMessage(id string) (*domain.Message, error) {
	zaplog.Logger.Info("модуль postgres: Чтение message", zaplog.LogString("id", id))

	message := &domain.Message{}
	query := "SELECT body FROM messenger WHERE id=$1"
	err := r.db.Get(message, query, id)
	if err != nil {
		zaplog.Logger.Error("модуль postgres: Ошибка чтения message", zaplog.LogString("id", id), zaplog.LogError(err))
		return nil, fmt.Errorf("модуль postgres: не удалось прочитать сообщение с id %s: %w", message.Id, err)
	}

	message.Id = id
	zaplog.Logger.Info("модуль postgres: Успешное Чтение message", zaplog.LogString("id", message.Id), zaplog.LogString("body", message.Body))

	return message, err
}

func (r *MessengerPostgresRepository) ReadMessages() ([]*domain.Message, error) {
	zaplog.Logger.Info("модуль postgres: Чтение messages")

	var slice []*domain.Message
	query := "SELECT id , body FROM messenger"
	err := r.db.Select(&slice, query)
	if err != nil {
		zaplog.Logger.Error("модуль postgres: Ошибка чтения messages")
		return nil, fmt.Errorf("модуль postgres: не удалось прочитать сообщения: %w", err)
	}
	zaplog.Logger.Info("модуль postgres: Messages прочитаны успешно", zaplog.LogString("number", strconv.Itoa(len(slice))))

	return slice, err
}
