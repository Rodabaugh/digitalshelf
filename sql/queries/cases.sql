-- name: CreateCase :one
INSERT INTO cases (id, created_at, updated_at, name, location_id)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2
)
RETURNING *;

-- name: GetCases :many
SELECT * FROM cases;

-- name: GetCasesByLocation :many
SELECT * FROM cases WHERE location_id = $1;

-- name: GetCaseByID :one
SELECT * FROM cases WHERE id = $1;

-- name: GetCaseLocation :one
SELECT locations.id, locations.name
FROM locations
JOIN cases ON locations.id = cases.location_id
WHERE cases.id = $1;