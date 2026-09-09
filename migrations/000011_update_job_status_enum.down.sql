-- Filename: 000011_update_job_status_enum.down.sql

BEGIN;

ALTER TYPE job_status ADD VALUE IF NOT EXISTS 'cancelled';

COMMIT;
