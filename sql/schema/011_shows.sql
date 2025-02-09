-- +goose Up
CREATE TABLE shows (id UUID PRIMARY KEY,
                        created_at TIMESTAMP NOT NULL,
                        updated_at TIMESTAMP NOT NULL,
                        title TEXT NOT NULL,
                        season INT NOT NULL,
                        genre TEXT NOT NULL,
                        actors TEXT NOT NULL,
                        writer TEXT NOT NULL,
                        director TEXT NOT NULL,
                        release_date DATE NOT NULL,
                        barcode TEXT NOT NULL,
                        shelf_id UUID NOT NULL REFERENCES shelves(id) ON DELETE CASCADE);

ALTER TABLE shows
ADD search tsvector
GENERATED ALWAYS AS (
    setweight(to_tsvector('english', title), 'A') || ' ' ||
    setweight(to_tsvector('simple', actors), 'B') || ' ' ||
    setweight(to_tsvector('english', genre), 'C') || ' ' ||
    setweight(to_tsvector('simple', writer), 'D') || ' ' ||
    setweight(to_tsvector('simple', director), 'D') :: tsvector
) STORED;

CREATE INDEX idx_shows_search ON shows USING GIN(search);

-- +goose Down
drop index idx_shows_search;
DROP TABLE shows;