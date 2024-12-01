-- name: CreateLocation :one
INSERT INTO locations (id, created_at, updated_at, name, owner_id)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2
)
RETURNING *;

-- name: GetLocations :many
SELECT * FROM locations;

-- name: GetLocationsByOwner :many
SELECT * FROM locations WHERE owner_id = $1;

-- name: GetLocationByID :one
SELECT * FROM locations WHERE id = $1;