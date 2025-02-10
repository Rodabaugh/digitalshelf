package main

import (
	"encoding/json"
	"net/http"

	"github.com/Rodabaugh/digitalshelf/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerMoviesUpdate(w http.ResponseWriter, r *http.Request) {
	// Get the movie from the database
	movieIDString := r.PathValue("movie_id")
	if movieIDString == "" {
		respondWithError(w, http.StatusBadRequest, "No movie id was provided", nil)
		return
	}

	movieID, err := uuid.Parse(movieIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid movie id format", err)
		return
	}

	var requestBody struct {
		ShelfID string `json:"shelf_id"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	shelfIDString := requestBody.ShelfID
	if shelfIDString == "" {
		respondWithError(w, http.StatusBadRequest, "No new shelf id was provided", nil)
		return
	}

	shelfID, err := uuid.Parse(shelfIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid shelf id format", err)
		return
	}

	// Validate user is authorized to modify movies at the location of requested movie.
	movieLocation, err := cfg.db.GetMovieLocation(r.Context(), movieID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get movie location", err)
		return
	}

	err = cfg.authorizeMember(movieLocation.ID, *r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to get movies at this location", err)
		return
	}

	// Validate user is authorized to modify movies at the location of new shelf.
	shelfLocation, err := cfg.db.GetShelfLocation(r.Context(), shelfID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get shelf location", err)
		return
	}

	err = cfg.authorizeMember(shelfLocation.ID, *r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to get movies at this location", err)
		return
	}

	movie, err := cfg.db.GetMovieByID(r.Context(), movieID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Movie not found", err)
		return
	}

	// Update the movie with the new shelf ID
	movie, err = cfg.db.UpdateMovie(r.Context(), database.UpdateMovieParams{
		ID:          movie.ID,
		Title:       movie.Title,
		Genre:       movie.Genre,
		Actors:      movie.Actors,
		Writer:      movie.Writer,
		Director:    movie.Director,
		ReleaseDate: movie.ReleaseDate,
		Barcode:     movie.Barcode,
		ShelfID:     shelfID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update movie", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Movie{
		ID:          movie.ID,
		Title:       movie.Title,
		Genre:       movie.Genre,
		Actors:      movie.Actors,
		Writer:      movie.Writer,
		Director:    movie.Director,
		Barcode:     movie.Barcode,
		ReleaseDate: movie.ReleaseDate,
		CreatedAt:   movie.CreatedAt,
		UpdatedAt:   movie.UpdatedAt,
		ShelfID:     movie.ShelfID,
	})
}
