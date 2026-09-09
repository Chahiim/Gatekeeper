-- Filename: 000008_create_images_table.up.sql

BEGIN;

CREATE TABLE IF NOT EXISTS images (
    id              bigserial    PRIMARY KEY,
    original_name   text         NOT NULL,
    stored_name     text         NOT NULL UNIQUE,
    media_type      text         NOT NULL,
    file_size       bigint       NOT NULL,
    created_at      timestamptz  NOT NULL DEFAULT now()
);

COMMIT;
