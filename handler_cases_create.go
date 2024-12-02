package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Rodabaugh/digitalshelf/internal/database"
	"github.com/google/uuid"
)

type Case struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	LocationID uuid.UUID `json:"location_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (cfg *apiConfig) handlerCasesCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name       string    `json:"name"`
		LocationID uuid.UUID `json:"location_id"`
	}
	type response struct {
		Case
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Was unable to decode parameters", err)
		return
	}

	item_case, err := cfg.db.CreateCase(r.Context(), database.CreateCaseParams{
		Name:       params.Name,
		LocationID: params.LocationID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create case", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Case: Case{
			ID:         item_case.ID,
			Name:       item_case.Name,
			LocationID: item_case.LocationID,
			CreatedAt:  item_case.CreatedAt,
			UpdatedAt:  item_case.UpdatedAt,
		},
	})
}
