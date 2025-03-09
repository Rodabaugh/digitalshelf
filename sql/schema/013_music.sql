-- +goose Up
CREATE TABLE music (id UUID PRIMARY KEY,
                        created_at TIMESTAMP NOT NULL,
                        updated_at TIMESTAMP NOT NULL,
                        title TEXT NOT NULL,
                        artist TEXT NOT NULL,
                        genre TEXT NOT NULL,
                        release_date DATE NOT NULL,
                        barcode TEXT NOT NULL,
                        format TEXT NOT NULL,
                        shelf_id UUID NOT NULL REFERENCES shelves(id) ON DELETE CASCADE);

ALTER TABLE music
ADD search tsvector
GENERATED ALWAYS AS (
    setweight(to_tsvector('english', title), 'A') || ' ' ||
    setweight(to_tsvector('simple', artist), 'B') || ' ' ||
    setweight(to_tsvector('english', genre), 'C') || ' ' ||
    setweight(to_tsvector('english', format), 'D') :: tsvector
) STORED;

CREATE INDEX idx_music_search ON music USING GIN(search);

-- +goose Down
DROP INDEX idx_music_search;
DROP TABLE music;