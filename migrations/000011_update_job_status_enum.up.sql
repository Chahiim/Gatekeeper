-- Filename: 000011_update_job_status_enum.up.sql

BEGIN;

ALTER TYPE job_status DROP VALUE IF EXISTS 'cancelled';

COMMIT;
