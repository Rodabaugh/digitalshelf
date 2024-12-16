-- +goose Up
CREATE TABLE movies (id UUID PRIMARY KEY,
                        created_at TIMESTAMP NOT NULL,
                        updated_at TIMESTAMP NOT NULL,
                        title TEXT NOT NULL,
                        genre TEXT NOT NULL,
                        actors TEXT NOT NULL,
                        writer TEXT NOT NULL,
                        director TEXT NOT NULL,
                        release_date DATE NOT NULL,
                        barcode TEXT NOT NULL,
                        shelf_id UUID NOT NULL REFERENCES shelves(id) ON DELETE CASCADE);

-- +goose Down
DROP TABLE movies;