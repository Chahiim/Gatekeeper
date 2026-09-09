package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type ImageJobStatus string

const (
	ImageJobStatusQueued     ImageJobStatus = "queued"
	ImageJobStatusProcessing ImageJobStatus = "processing"
	ImageJobStatusCompleted  ImageJobStatus = "completed"
	ImageJobStatusFailed     ImageJobStatus = "failed"
)

type ImageJob struct {
	ID           int64         `json:"id"`
	ImageID      int64         `json:"image_id"`
	Status       ImageJobStatus `json:"status"`
	ErrorMessage *string       `json:"error_message,omitempty"`
	QueuedAt     time.Time     `json:"queued_at"`
	StartedAt    *time.Time    `json:"started_at,omitempty"`
	CompletedAt  *time.Time    `json:"completed_at,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
}

type ImageJobModel struct {
	DB *sql.DB
}

func (m ImageJobModel) Insert(job *ImageJob) error {
	query := `
		INSERT INTO image_jobs (image_id, status)
		VALUES ($1, $2)
		RETURNING id, queued_at, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query,
		job.ImageID,
		job.Status,
	).Scan(&job.ID, &job.QueuedAt, &job.CreatedAt)
}

func (m ImageJobModel) Get(id int64) (*ImageJob, error) {
	query := `
		SELECT id, image_id, status, error_message, queued_at, started_at, completed_at, created_at
		FROM image_jobs
		WHERE id = $1`

	var job ImageJob

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&job.ID,
		&job.ImageID,
		&job.Status,
		&job.ErrorMessage,
		&job.QueuedAt,
		&job.StartedAt,
		&job.CompletedAt,
		&job.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &job, nil
}

func (m ImageJobModel) GetByImageID(imageID int64) (*ImageJob, error) {
	query := `
		SELECT id, image_id, status, error_message, queued_at, started_at, completed_at, created_at
		FROM image_jobs
		WHERE image_id = $1
		ORDER BY created_at DESC
		LIMIT 1`

	var job ImageJob

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, imageID).Scan(
		&job.ID,
		&job.ImageID,
		&job.Status,
		&job.ErrorMessage,
		&job.QueuedAt,
		&job.StartedAt,
		&job.CompletedAt,
		&job.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &job, nil
}
