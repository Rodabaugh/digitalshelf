-- +goose Up
CREATE TABLE shelves (id UUID PRIMARY KEY,
                        created_at TIMESTAMP NOT NULL,
                        updated_at TIMESTAMP NOT NULL,
                        name TEXT NOT NULL,
                        case_id UUID NOT NULL REFERENCES cases(id) ON DELETE CASCADE);

-- +goose Down
DROP TABLE shelves;