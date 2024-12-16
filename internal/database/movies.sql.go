// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: movies.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createMovie = `-- name: CreateMovie :one
INSERT INTO movies (id, created_at, updated_at, title, genre, actors, writer, director, release_date, barcode, shelf_id)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id, created_at, updated_at, title, genre, actors, writer, director, release_date, barcode, shelf_id
`

type CreateMovieParams struct {
	Title       string
	Genre       string
	Actors      string
	Writer      string
	Director    string
	ReleaseDate time.Time
	Barcode     string
	ShelfID     uuid.UUID
}

func (q *Queries) CreateMovie(ctx context.Context, arg CreateMovieParams) (Movie, error) {
	row := q.db.QueryRowContext(ctx, createMovie,
		arg.Title,
		arg.Genre,
		arg.Actors,
		arg.Writer,
		arg.Director,
		arg.ReleaseDate,
		arg.Barcode,
		arg.ShelfID,
	)
	var i Movie
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Genre,
		&i.Actors,
		&i.Writer,
		&i.Director,
		&i.ReleaseDate,
		&i.Barcode,
		&i.ShelfID,
	)
	return i, err
}

const getMovieByID = `-- name: GetMovieByID :one
SELECT id, created_at, updated_at, title, genre, actors, writer, director, release_date, barcode, shelf_id FROM movies WHERE id = $1
`

func (q *Queries) GetMovieByID(ctx context.Context, id uuid.UUID) (Movie, error) {
	row := q.db.QueryRowContext(ctx, getMovieByID, id)
	var i Movie
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Genre,
		&i.Actors,
		&i.Writer,
		&i.Director,
		&i.ReleaseDate,
		&i.Barcode,
		&i.ShelfID,
	)
	return i, err
}

const getMovies = `-- name: GetMovies :many
SELECT id, created_at, updated_at, title, genre, actors, writer, director, release_date, barcode, shelf_id FROM movies
`

func (q *Queries) GetMovies(ctx context.Context) ([]Movie, error) {
	rows, err := q.db.QueryContext(ctx, getMovies)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Movie
	for rows.Next() {
		var i Movie
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Genre,
			&i.Actors,
			&i.Writer,
			&i.Director,
			&i.ReleaseDate,
			&i.Barcode,
			&i.ShelfID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMoviesByShelf = `-- name: GetMoviesByShelf :many
SELECT id, created_at, updated_at, title, genre, actors, writer, director, release_date, barcode, shelf_id FROM movies WHERE shelf_id = $1
`

func (q *Queries) GetMoviesByShelf(ctx context.Context, shelfID uuid.UUID) ([]Movie, error) {
	rows, err := q.db.QueryContext(ctx, getMoviesByShelf, shelfID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Movie
	for rows.Next() {
		var i Movie
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Genre,
			&i.Actors,
			&i.Writer,
			&i.Director,
			&i.ReleaseDate,
			&i.Barcode,
			&i.ShelfID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
