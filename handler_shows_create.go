package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Rodabaugh/digitalshelf/internal/database"
	"github.com/google/uuid"
)

type Show struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Season      string    `json:"season"`
	Genre       string    `json:"genre"`
	Actors      string    `json:"actors"`
	Writer      string    `json:"writer"`
	Director    string    `json:"director"`
	Barcode     string    `json:"barcode"`
	Format      string    `json:"format"`
	ShelfID     uuid.UUID `json:"shelf_id"`
	ReleaseDate time.Time `json:"release_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (cfg *apiConfig) handlerShowCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title       string    `json:"title"`
		Season      string    `json:"season"`
		Genre       string    `json:"genre"`
		Actors      string    `json:"actors"`
		Writer      string    `json:"writer"`
		Director    string    `json:"director"`
		Barcode     string    `json:"barcode"`
		Format      string    `json:"format"`
		ShelfID     uuid.UUID `json:"shelf_id"`
		ReleaseDate time.Time `json:"release_date"`
	}

	type response struct {
		Show
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
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to create shows in this location", err)
		return
	}

	show, err := cfg.db.CreateShow(r.Context(), database.CreateShowParams{
		Title:       params.Title,
		Season:      params.Season,
		Genre:       params.Genre,
		Actors:      params.Actors,
		Writer:      params.Writer,
		Director:    params.Director,
		Barcode:     params.Barcode,
		Format:      params.Format,
		ShelfID:     params.ShelfID,
		ReleaseDate: params.ReleaseDate,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create show", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Show: Show{
			ID:          show.ID,
			Title:       show.Title,
			Season:      show.Season,
			Genre:       show.Genre,
			Actors:      show.Actors,
			Writer:      show.Writer,
			Director:    show.Director,
			Barcode:     show.Barcode,
			Format:      show.Format,
			ShelfID:     show.ShelfID,
			ReleaseDate: show.ReleaseDate,
			CreatedAt:   show.CreatedAt,
			UpdatedAt:   show.UpdatedAt,
		},
	})
}
