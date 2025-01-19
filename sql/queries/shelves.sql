-- name: CreateShelf :one
INSERT INTO shelves (id, created_at, updated_at, name, case_id)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2
)
RETURNING *;

-- name: GetShelves :many
SELECT * FROM shelves;

-- name: GetShelvesByCase :many
SELECT * FROM shelves WHERE case_id = $1;

-- name: GetShelfByID :one
SELECT * FROM shelves WHERE id = $1;

-- name: GetShelfLocation :one
SELECT locations.id, locations.name
FROM locations
JOIN cases ON locations.id = cases.location_id
JOIN shelves ON cases.id = shelves.case_id
WHERE shelves.id = $1;