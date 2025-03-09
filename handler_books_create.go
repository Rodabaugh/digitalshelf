package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Rodabaugh/digitalshelf/internal/database"
	"github.com/google/uuid"
)

type Book struct {
	ID              uuid.UUID `json:"id"`
	Title           string    `json:"title"`
	Author          string    `json:"author"`
	Genre           string    `json:"genre"`
	Barcode         string    `json:"barcode"`
	ShelfID         uuid.UUID `json:"shelf_id"`
	PublicationDate time.Time `json:"publication_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (cfg *apiConfig) handlerBookCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title           string    `json:"title"`
		Author          string    `json:"author"`
		Genre           string    `json:"genre"`
		Barcode         string    `json:"barcode"`
		ShelfID         uuid.UUID `json:"shelf_id"`
		PublicationDate time.Time `json:"publication_date"`
	}

	type response struct {
		Book
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Was unable to decode parameters", err)
		return
	}

	shelfLocation, err := cfg.db.GetShelfLocation(r.Context(), params.ShelfID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get shelf location", err)
		return
	}

	err = cfg.authorizeMember(shelfLocation.ID, *r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to create books in this location", err)
		return
	}

	book, err := cfg.db.CreateBook(r.Context(), database.CreateBookParams{
		Title:           params.Title,
		Author:          params.Author,
		Genre:           params.Genre,
		Barcode:         params.Barcode,
		ShelfID:         params.ShelfID,
		PublicationDate: params.PublicationDate,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create book", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Book: Book{
			ID:              book.ID,
			Title:           book.Title,
			Author:          book.Author,
			Genre:           book.Genre,
			Barcode:         book.Barcode,
			ShelfID:         book.ShelfID,
			PublicationDate: book.PublicationDate,
			CreatedAt:       book.CreatedAt,
			UpdatedAt:       book.UpdatedAt,
		},
	})
}
