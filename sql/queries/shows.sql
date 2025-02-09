-- name: CreateShow :one
INSERT INTO shows (id, created_at, updated_at, title, season, genre, actors, writer, director, release_date, barcode, shelf_id)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetShows :many
SELECT * FROM shows;

-- name: GetShowsByShelf :many
SELECT * FROM shows WHERE shelf_id = $1;

-- name: GetShowByID :one
SELECT * FROM shows WHERE id = $1;

-- name: GetShowByBarcode :one
SELECT * FROM shows WHERE barcode = $1;

-- name: GetShowsByLocation :many
SELECT * FROM shows
INNER JOIN shelves
ON shows.shelf_id = shelves.id
INNER JOIN cases
ON shelves.case_id = cases.id
INNER JOIN locations
ON cases.location_id = locations.id
WHERE locations.id = $1;

-- name: GetShowLocation :one
SELECT locations.id, locations.name
FROM locations
JOIN cases ON locations.id = cases.location_id
JOIN shelves ON cases.id = shelves.case_id
JOIN shows ON shelves.id = shows.shelf_id
WHERE shows.id = $1;

-- name: SearchShows :many
SELECT shows.id, shows.created_at, shows.updated_at, title, season, genre, actors, writer, director, release_date, barcode, shelf_id,
    CAST(
        ts_rank(search, websearch_to_tsquery('english', $1)) + 
        ts_rank(search, websearch_to_tsquery('simple', $1)) AS float8
    ) AS rank
FROM shows
INNER JOIN shelves
ON shows.shelf_id = shelves.id
INNER JOIN cases
ON shelves.case_id = cases.id
INNER JOIN locations
ON cases.location_id = locations.id
WHERE search @@ websearch_to_tsquery('english', $1)
OR search @@ websearch_to_tsquery('simple', $1)
AND locations.id = $2
ORDER BY rank DESC;