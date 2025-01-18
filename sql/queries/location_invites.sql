-- name: AddLocationInvite :one
INSERT INTO location_invites (location_id, user_id, invited_at)
VALUES (
    $1, $2, NOW()
)
RETURNING *;

-- name: GetLocationInvites :many
SELECT location_invites.location_id, users.id, users.name, users.email, location_invites.invited_at
FROM users
INNER JOIN location_invites
ON users.ID = location_invites.user_id
WHERE location_invites.location_id = $1;

-- name: GetUserInvites :many
SELECT location_invites.user_id, locations.id, locations.name, locations.owner_id, location_invites.invited_at
FROM locations
INNER JOIN location_invites
ON locations.ID = location_invites.location_id
WHERE location_invites.user_id = $1;

-- name: RemoveLocationInvite :exec
DELETE FROM location_invites
WHERE location_id = $1 AND user_id = $2;