-- +goose Up
ALTER TABLE movies
ADD format TEXT NOT NULL DEFAULT 'Not Specified'; 

-- +goose Down
ALTER TABLE movies
DROP COLUMN format;