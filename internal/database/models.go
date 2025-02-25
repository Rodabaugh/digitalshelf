// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Case struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Name       string
	LocationID uuid.UUID
}

type Location struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	OwnerID   uuid.UUID
}

type LocationInvite struct {
	LocationID uuid.UUID
	UserID     uuid.UUID
	InvitedAt  time.Time
}

type LocationUser struct {
	LocationID uuid.UUID
	UserID     uuid.UUID
	JoinedAt   time.Time
}

type Movie struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Genre       string
	Actors      string
	Writer      string
	Director    string
	ReleaseDate time.Time
	Barcode     string
	ShelfID     uuid.UUID
	Search      interface{}
}

type RefreshToken struct {
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	ExpiresAt time.Time
	RevokedAt sql.NullTime
}

type Shelf struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	CaseID    uuid.UUID
}

type Show struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Season      int32
	Genre       string
	Actors      string
	Writer      string
	Director    string
	ReleaseDate time.Time
	Barcode     string
	ShelfID     uuid.UUID
	Search      interface{}
}

type User struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Name           string
	Email          string
	HashedPassword string
}
