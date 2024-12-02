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