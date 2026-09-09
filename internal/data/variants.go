package data

import (
	"context"
	"database/sql"
	"time"
)

type Variant struct {
	ID         int64     `json:"id"`
	ImageID    int64     `json:"image_id"`
	Name       string    `json:"name"`
	StoredName string    `json:"stored_name"`
	Width      int       `json:"width"`
	Height     int       `json:"height"`
	FileSize   int64     `json:"file_size"`
	CreatedAt  time.Time `json:"created_at"`
}

type VariantModel struct {
	DB *sql.DB
}

func (m VariantModel) Insert(variant *Variant) error {
	query := `
		INSERT INTO variants (image_id, name, stored_name, width, height, file_size)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query,
		variant.ImageID,
		variant.Name,
		variant.StoredName,
		variant.Width,
		variant.Height,
		variant.FileSize,
	).Scan(&variant.ID, &variant.CreatedAt)
}

func (m VariantModel) GetByImageID(imageID int64) ([]Variant, error) {
	query := `
		SELECT id, image_id, name, stored_name, width, height, file_size, created_at
		FROM variants
		WHERE image_id = $1
		ORDER BY name`

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, imageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var variants []Variant
	for rows.Next() {
		var v Variant
		err := rows.Scan(
			&v.ID,
			&v.ImageID,
			&v.Name,
			&v.StoredName,
			&v.Width,
			&v.Height,
			&v.FileSize,
			&v.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		variants = append(variants, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return variants, nil
}
