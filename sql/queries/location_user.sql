-- name: AddLocationMember :one
INSERT INTO location_user (location_id, user_id, joined_at)
VALUES (
    $1, $2, NOW()
)
RETURNING *;

-- name: GetLocationMembers :many
SELECT location_user.location_id, users.id, users.name, Users.email, location_user.joined_at
FROM users
INNER JOIN location_user
ON users.ID = location_user.user_id
WHERE location_user.location_id = $1;

-- name: GetUserLocations :many
SELECT location_user.user_id, locations.id, locations.name, locations.owner_id, location_user.joined_at
FROM locations
INNER JOIN location_user
ON locations.ID = location_user.location_id
WHERE location_user.user_id = $1;

-- name: RemoveLocationMember :exec
DELETE FROM location_user
WHERE location_id = $1 AND user_id = $2;