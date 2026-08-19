package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type ConsumerStatus string

const (
	ConsumerStatusActive    ConsumerStatus = "active"
	ConsumerStatusSuspended ConsumerStatus = "suspended"
	ConsumerStatusTerminated ConsumerStatus = "terminated"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Consumer struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Status    ConsumerStatus `json:"status"`
	Version   int64          `json:"version"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type ConsumerModel struct {
	DB *sql.DB
}

func (m ConsumerModel) Insert(consumer *Consumer) error {
	query := `
		INSERT INTO consumers (name, email, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query,
		consumer.Name,
		consumer.Email,
		consumer.Status,
	).Scan(&consumer.ID, &consumer.CreatedAt, &consumer.UpdatedAt)
}

func (m ConsumerModel) Get(id string) (*Consumer, error) {
	query := `
		SELECT id, name, email, status, version, created_at, updated_at
		FROM consumers
		WHERE id = $1`

	var consumer Consumer

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&consumer.ID,
		&consumer.Name,
		&consumer.Email,
		&consumer.Status,
		&consumer.Version,
		&consumer.CreatedAt,
		&consumer.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &consumer, nil
}
