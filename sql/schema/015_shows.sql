-- +goose Up
ALTER TABLE shows
ADD COLUMN format TEXT NOT NULL DEFAULT 'Not Specified';

ALTER TABLE shows
ALTER COLUMN season TYPE TEXT;

-- +goose Down
ALTER TABLE movies
DROP COLUMN format;
-- This migration does not support a down migration of the season column type. Seasons shall hereafter be stored as text, to prevent data loss and to better support box-sets.