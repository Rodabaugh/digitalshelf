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

-- name: GetMovieByBarcode :one
SELECT * FROM movies WHERE barcode = $1;

-- name: GetMoviesByLocation :many
SELECT * FROM movies
INNER JOIN shelves
ON movies.shelf_id = shelves.id
INNER JOIN cases
ON shelves.case_id = cases.id
INNER JOIN locations
ON cases.location_id = locations.id
WHERE locations.id = $1;

-- name: GetMovieLocation :one
SELECT locations.id, locations.name
FROM locations
JOIN cases ON locations.id = cases.location_id
JOIN shelves ON cases.id = shelves.case_id
JOIN movies ON shelves.id = movies.shelf_id
WHERE movies.id = $1;