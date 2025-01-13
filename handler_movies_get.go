package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerMoviesGet(w http.ResponseWriter, r *http.Request) {
	dbMovies, err := cfg.db.GetMovies(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get movies from database", err)
		return
	}

	movies := []Movie{}

	for _, dbMovie := range dbMovies {
		movies = append(movies, Movie{
			ID:          dbMovie.ID,
			Title:       dbMovie.Title,
			Genre:       dbMovie.Genre,
			Actors:      dbMovie.Actors,
			Writer:      dbMovie.Writer,
			Director:    dbMovie.Director,
			Barcode:     dbMovie.Barcode,
			ReleaseDate: dbMovie.ReleaseDate,
			CreatedAt:   dbMovie.CreatedAt,
			UpdatedAt:   dbMovie.UpdatedAt,
			ShelfID:     dbMovie.ShelfID,
		})
	}

	respondWithJSON(w, http.StatusOK, movies)
}

func (cfg *apiConfig) handlerMoviesGetByShelf(w http.ResponseWriter, r *http.Request) {
	shelfIDString := r.PathValue("shelf_id")
	if shelfIDString == "" {
		respondWithError(w, http.StatusBadRequest, "No shelf id was provided", fmt.Errorf("no shelf_id was provided"))
		return
	}

	shelfID, err := uuid.Parse(shelfIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid shelf ID", err)
		return
	}

	dbMovies, err := cfg.db.GetMoviesByShelf(r.Context(), shelfID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No movies found for that shelf", err)
		return
	}

	movies := []Movie{}

	for _, dbMovie := range dbMovies {
		movies = append(movies, Movie{
			ID:          dbMovie.ID,
			Title:       dbMovie.Title,
			Genre:       dbMovie.Genre,
			Actors:      dbMovie.Actors,
			Writer:      dbMovie.Writer,
			Director:    dbMovie.Director,
			Barcode:     dbMovie.Barcode,
			ReleaseDate: dbMovie.ReleaseDate,
			CreatedAt:   dbMovie.CreatedAt,
			UpdatedAt:   dbMovie.UpdatedAt,
			ShelfID:     dbMovie.ShelfID,
		})
	}

	respondWithJSON(w, http.StatusOK, movies)
}

func (cfg *apiConfig) handlerMovieGetByID(w http.ResponseWriter, r *http.Request) {
	movieIDString := r.PathValue("movie_id")
	if movieIDString == "" {
		respondWithError(w, http.StatusBadRequest, "No movie id was provided", fmt.Errorf("no movie id was provided"))
		return
	}

	movieID, err := uuid.Parse(movieIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid movie ID", err)
		return
	}

	dbMovie, err := cfg.db.GetMovieByID(r.Context(), movieID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Movie not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Movie{
		ID:          dbMovie.ID,
		Title:       dbMovie.Title,
		Genre:       dbMovie.Genre,
		Actors:      dbMovie.Actors,
		Writer:      dbMovie.Writer,
		Director:    dbMovie.Director,
		Barcode:     dbMovie.Barcode,
		ReleaseDate: dbMovie.ReleaseDate,
		CreatedAt:   dbMovie.CreatedAt,
		UpdatedAt:   dbMovie.UpdatedAt,
		ShelfID:     dbMovie.ShelfID,
	})
}

func (cfg *apiConfig) handlerGetMovieByBarcode(w http.ResponseWriter, r *http.Request) {
	barcode := r.PathValue("barcode")
	fmt.Println(barcode)
	if barcode == "" {
		respondWithError(w, http.StatusBadRequest, "No barcode was provided", fmt.Errorf("no barcode was provided"))
		return
	}

	dbMovie, err := cfg.db.GetMovieByBarcode(r.Context(), barcode)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Movie not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Movie{
		ID:          dbMovie.ID,
		Title:       dbMovie.Title,
		Genre:       dbMovie.Genre,
		Actors:      dbMovie.Actors,
		Writer:      dbMovie.Writer,
		Director:    dbMovie.Director,
		Barcode:     dbMovie.Barcode,
		ReleaseDate: dbMovie.ReleaseDate,
		CreatedAt:   dbMovie.CreatedAt,
		UpdatedAt:   dbMovie.UpdatedAt,
		ShelfID:     dbMovie.ShelfID,
	})
}
