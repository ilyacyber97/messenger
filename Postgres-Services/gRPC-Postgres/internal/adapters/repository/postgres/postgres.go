package postgres

import (
	"fmt"
	"grpc-postgres/domain"
	"grpc-postgres/internal/core/ports"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type PostgresRepository struct {
	db  *sqlx.DB
	log ports.LoggerPort
}

func NewPostgresRepository(DSN string, log ports.LoggerPort) (*PostgresRepository, error) {
	db, err := sqlx.Connect("postgres", DSN)
	if err != nil {
		return nil, fmt.Errorf("postgres: ошибка подключения к PostgreSQL: %w", err)
	}

	return &PostgresRepository{db: db, log: log}, err
}

func (r *PostgresRepository) SaveMessage(message *domain.Message) error {
	r.log.Info("postgres: Сохранение message", zap.String("id", message.Id),
		zap.String("body", message.Body))

	query := "INSERT INTO messenger(id ,body) VALUES(:id, :body)"
	_, err := r.db.NamedExec(query, &message)
	if err != nil {
		r.log.Error("postgres: Ошибка сохранения message", zap.String("id", message.Id), zap.Error(err))
		return fmt.Errorf("postgres: не удалось сохранить сообщение с id %s: %w", message.Id, err)

	}

	r.log.Info("postgres: Успешная Запись message", zap.String("id", message.Id), zap.String("body", message.Body))

	return err
}

func (r *PostgresRepository) ReadMessage(id string) (*domain.Message, error) {
	r.log.Info("postgres: Чтение message", zap.String("id", id))

	message := &domain.Message{}
	query := "SELECT body FROM messenger WHERE id=$1"
	err := r.db.Get(message, query, id)
	if err != nil {
		r.log.Error("postgres: Ошибка чтения message", zap.String("id", id), zap.Error(err))
		return nil, fmt.Errorf("postgres: не удалось прочитать сообщение с id %s: %w", message.Id, err)
	}

	message.Id = id
	r.log.Info("postgres: Успешное Чтение message", zap.String("id", message.Id), zap.String("body", message.Body))

	return message, err
}

func (r *PostgresRepository) ReadMessages() ([]*domain.Message, error) {
	r.log.Info("postgres: Чтение messages")

	var slice []*domain.Message
	query := "SELECT id , body FROM messenger"
	err := r.db.Select(&slice, query)
	if err != nil {
		r.log.Error("postgres: Ошибка чтения messages")
		return nil, fmt.Errorf("postgres: не удалось прочитать сообщения: %w", err)
	}
	r.log.Info("postgres: Messages прочитаны успешно", zap.String("number", strconv.Itoa(len(slice))))

	return slice, err
}
