-- +goose Up
CREATE TABLE cases (id UUID PRIMARY KEY,
                        created_at TIMESTAMP NOT NULL,
                        updated_at TIMESTAMP NOT NULL,
                        name TEXT NOT NULL,
                        location_id UUID NOT NULL REFERENCES locations(id) ON DELETE CASCADE);

-- +goose Down
DROP TABLE cases;