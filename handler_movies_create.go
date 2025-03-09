package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Rodabaugh/digitalshelf/internal/database"
	"github.com/google/uuid"
)

type Movie struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
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

func (cfg *apiConfig) handlerMovieCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title       string    `json:"title"`
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
		Movie
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
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to create movies in this location", err)
		return
	}

	movie, err := cfg.db.CreateMovie(r.Context(), database.CreateMovieParams{
		Title:       params.Title,
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
		respondWithError(w, http.StatusInternalServerError, "Unable to create movie", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Movie: Movie{
			ID:          movie.ID,
			Title:       movie.Title,
			Genre:       movie.Genre,
			Actors:      movie.Actors,
			Writer:      movie.Writer,
			Director:    movie.Director,
			Barcode:     movie.Barcode,
			Format:      movie.Format,
			ShelfID:     movie.ShelfID,
			ReleaseDate: movie.ReleaseDate,
			CreatedAt:   movie.CreatedAt,
			UpdatedAt:   movie.UpdatedAt,
		},
	})
}
