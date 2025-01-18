-- +goose Up
CREATE TABLE location_invites (location_id UUID NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
                        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                        invited_at TIMESTAMP NOT NULL,
                        UNIQUE(location_id, user_id));

-- +goose Down
DROP TABLE location_invites;