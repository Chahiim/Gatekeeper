-- Filename: 000010_create_variants_table.up.sql

BEGIN;

CREATE TABLE IF NOT EXISTS variants (
    id              bigserial    PRIMARY KEY,
    image_id        bigint       NOT NULL REFERENCES images(id) ON DELETE CASCADE,
    name            text         NOT NULL,
    stored_name     text         NOT NULL UNIQUE,
    width           integer      NOT NULL,
    height          integer      NOT NULL,
    file_size       bigint       NOT NULL,
    created_at      timestamptz  NOT NULL DEFAULT now(),
    UNIQUE(image_id, name)
);

COMMIT;
