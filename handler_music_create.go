package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Rodabaugh/digitalshelf/internal/database"
	"github.com/google/uuid"
)

type Music struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Artist      string    `json:"artist"`
	Genre       string    `json:"genre"`
	Barcode     string    `json:"barcode"`
	Format      string    `json:"format"`
	ShelfID     uuid.UUID `json:"shelf_id"`
	ReleaseDAte time.Time `json:"release_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (cfg *apiConfig) handlerMusicCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title       string    `json:"title"`
		Artist      string    `json:"artist"`
		Genre       string    `json:"genre"`
		Barcode     string    `json:"barcode"`
		Format      string    `json:"format"`
		ShelfID     uuid.UUID `json:"shelf_id"`
		ReleaseDate time.Time `json:"release_date"`
	}

	type response struct {
		Music
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
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to create music in this location", err)
		return
	}

	music, err := cfg.db.CreateMusic(r.Context(), database.CreateMusicParams{
		Title:       params.Title,
		Artist:      params.Artist,
		Genre:       params.Genre,
		Barcode:     params.Barcode,
		Format:      params.Format,
		ShelfID:     params.ShelfID,
		ReleaseDate: params.ReleaseDate,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create music", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Music: Music{
			ID:          music.ID,
			Title:       music.Title,
			Artist:      music.Artist,
			Genre:       music.Genre,
			Barcode:     music.Barcode,
			Format:      music.Format,
			ShelfID:     music.ShelfID,
			ReleaseDAte: music.ReleaseDate,
			CreatedAt:   music.CreatedAt,
			UpdatedAt:   music.UpdatedAt,
		},
	})
}
