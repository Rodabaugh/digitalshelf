package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Rodabaugh/digitalshelf/internal/database"
	"github.com/google/uuid"
)

type Location struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	OwnerID   uuid.UUID `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (cfg *apiConfig) handlerLocationsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name    string    `json:"name"`
		OwnerID uuid.UUID `json:"owner_id"`
	}
	type response struct {
		Location
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Was unable to decode parameters", err)
		return
	}

	location, err := cfg.db.CreateLocation(r.Context(), database.CreateLocationParams{
		Name:    params.Name,
		OwnerID: params.OwnerID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create location", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Location: Location{
			ID:        location.ID,
			Name:      location.Name,
			OwnerID:   location.OwnerID,
			CreatedAt: location.CreatedAt,
			UpdatedAt: location.UpdatedAt,
		},
	})
}
