-- +goose Up
CREATE TABLE books (id UUID PRIMARY KEY,
                        created_at TIMESTAMP NOT NULL,
                        updated_at TIMESTAMP NOT NULL,
                        title TEXT NOT NULL,
                        author TEXT NOT NULL,
                        genre TEXT NOT NULL,
                        publication_date DATE NOT NULL,
                        barcode TEXT NOT NULL,
                        shelf_id UUID NOT NULL REFERENCES shelves(id) ON DELETE CASCADE);

ALTER TABLE books
ADD search tsvector
GENERATED ALWAYS AS (
    setweight(to_tsvector('english', title), 'A') || ' ' ||
    setweight(to_tsvector('simple', author), 'B') || ' ' ||
    setweight(to_tsvector('english', genre), 'C') :: tsvector
) STORED;

CREATE INDEX idx_books_search ON books USING GIN(search);

-- +goose Down
DROP INDEX idx_books_search;
DROP TABLE books;