package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Rodabaugh/digitalshelf/internal/database"
	"github.com/google/uuid"
)

type Shelf struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CaseID    uuid.UUID `json:"case_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (cfg *apiConfig) handlerShelfCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name   string    `json:"name"`
		CaseID uuid.UUID `json:"case_id"`
	}
	type response struct {
		Shelf
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Was unable to decode parameters", err)
		return
	}

	shelf, err := cfg.db.CreateShelf(r.Context(), database.CreateShelfParams{
		Name:   params.Name,
		CaseID: params.CaseID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create shelf", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Shelf: Shelf{
			ID:        shelf.ID,
			Name:      shelf.Name,
			CaseID:    shelf.CaseID,
			CreatedAt: shelf.CreatedAt,
			UpdatedAt: shelf.UpdatedAt,
		},
	})
}
