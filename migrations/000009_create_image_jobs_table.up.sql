-- Filename: 000009_create_image_jobs_table.up.sql

BEGIN;

CREATE TABLE IF NOT EXISTS image_jobs (
    id              bigserial    PRIMARY KEY,
    image_id        bigint       NOT NULL REFERENCES images(id) ON DELETE CASCADE,
    status          job_status   NOT NULL DEFAULT 'queued',
    error_message   text,
    queued_at       timestamptz  NOT NULL DEFAULT now(),
    started_at      timestamptz,
    completed_at    timestamptz,
    created_at      timestamptz  NOT NULL DEFAULT now()
);

COMMIT;
