-- name: CreateMusic :one
INSERT INTO music (id, created_at, updated_at, title, artist, genre, release_date, barcode, format, shelf_id)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetMusic :many
SELECT * FROM music;

-- name: GetMusicByShelf :many
SELECT * FROM music WHERE shelf_id = $1;

-- name: GetMusicByID :one
SELECT * FROM music WHERE id = $1;

-- name: GetMusicByBarcode :one
SELECT * FROM music WHERE barcode = $1;

-- name: GetMusicByLocation :many
SELECT music.id, music.created_at, music.updated_at, title, artist, genre, release_date, barcode, format, shelf_id
FROM music
INNER JOIN shelves
ON music.shelf_id = shelves.id
INNER JOIN cases
ON shelves.case_id = cases.id
INNER JOIN locations
ON cases.location_id = locations.id
WHERE locations.id = $1;

-- name: GetMusicLocation :one
SELECT locations.id, locations.name
FROM locations
JOIN cases ON locations.id = cases.location_id
JOIN shelves ON cases.id = shelves.case_id
JOIN music ON shelves.id = music.shelf_id
WHERE music.id = $1;

-- name: SearchMusic :many
SELECT music.id, music.created_at, music.updated_at, title, artist, genre, release_date, barcode, format, shelf_id,
    CAST(
        ts_rank(search, websearch_to_tsquery('english', $1)) + 
        ts_rank(search, websearch_to_tsquery('simple', $1)) AS float8
    ) AS rank
FROM music
INNER JOIN shelves
ON music.shelf_id = shelves.id
INNER JOIN cases
ON shelves.case_id = cases.id
INNER JOIN locations
ON cases.location_id = locations.id
WHERE search @@ websearch_to_tsquery('english', $1)
OR search @@ websearch_to_tsquery('simple', $1)
AND locations.id = $2
ORDER BY rank DESC;