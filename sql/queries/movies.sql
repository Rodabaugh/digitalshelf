-- name: CreateMovie :one
INSERT INTO movies (id, created_at, updated_at, title, genre, actors, writer, director, release_date, barcode, shelf_id)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetMovies :many
SELECT * FROM movies;

-- name: GetMoviesByShelf :many
SELECT * FROM movies WHERE shelf_id = $1;

-- name: GetMovieByID :one
SELECT * FROM movies WHERE id = $1;