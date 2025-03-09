-- name: CreateBook :one
INSERT INTO books (id, created_at, updated_at, title, author, genre, publication_date, barcode, shelf_id)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetBooks :many
SELECT * FROM books;

-- name: GetBooksByShelf :many
SELECT * FROM books WHERE shelf_id = $1;

-- name: GetBookByID :one
SELECT * FROM books WHERE id = $1;

-- name: GetBookByBarcode :one
SELECT * FROM books WHERE barcode = $1;

-- name: GetBooksByLocation :many
SELECT * FROM books
INNER JOIN shelves
ON books.shelf_id = shelves.id
INNER JOIN cases
ON shelves.case_id = cases.id
INNER JOIN locations
ON cases.location_id = locations.id
WHERE locations.id = $1;

-- name: GetBookLocation :one
SELECT locations.id, locations.name
FROM locations
JOIN cases ON locations.id = cases.location_id
JOIN shelves ON cases.id = shelves.case_id
JOIN books ON shelves.id = books.shelf_id
WHERE books.id = $1;

-- name: SearchBooks :many
SELECT books.id, books.created_at, books.updated_at, title, author, genre, publication_date, barcode, shelf_id,
    CAST(
        ts_rank(search, websearch_to_tsquery('english', $1)) + 
        ts_rank(search, websearch_to_tsquery('simple', $1)) AS float8
    ) AS rank
FROM books
INNER JOIN shelves
ON books.shelf_id = shelves.id
INNER JOIN cases
ON shelves.case_id = cases.id
INNER JOIN locations
ON cases.location_id = locations.id
WHERE search @@ websearch_to_tsquery('english', $1)
OR search @@ websearch_to_tsquery('simple', $1)
AND locations.id = $2
ORDER BY rank DESC;