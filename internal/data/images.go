package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Image struct {
	ID           int64     `json:"id"`
	OriginalName string    `json:"original_name"`
	StoredName   string    `json:"stored_name"`
	MediaType    string    `json:"media_type"`
	FileSize     int64     `json:"file_size"`
	CreatedAt    time.Time `json:"created_at"`
}

type ImageModel struct {
	DB *sql.DB
}

func (m ImageModel) Insert(image *Image) error {
	query := `
		INSERT INTO images (original_name, stored_name, media_type, file_size)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query,
		image.OriginalName,
		image.StoredName,
		image.MediaType,
		image.FileSize,
	).Scan(&image.ID, &image.CreatedAt)
}

func (m ImageModel) Get(id int64) (*Image, error) {
	query := `
		SELECT id, original_name, stored_name, media_type, file_size, created_at
		FROM images
		WHERE id = $1`

	var image Image

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&image.ID,
		&image.OriginalName,
		&image.StoredName,
		&image.MediaType,
		&image.FileSize,
		&image.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &image, nil
}

func (m ImageModel) GetStoredName(id int64) (string, error) {
	query := `SELECT stored_name FROM images WHERE id = $1`

	var storedName string
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&storedName)
	if err != nil {
		return "", fmt.Errorf("could not get stored name: %w", err)
	}

	return storedName, nil
}
