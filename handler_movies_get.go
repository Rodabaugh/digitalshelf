package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rodabaugh/digitalshelf/internal/database"
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
			Format:      dbMovie.Format,
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

	// Validate user is authorized to get movies at the location of shelf.
	shelfLocation, err := cfg.db.GetShelfLocation(r.Context(), shelfID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get shelf location", err)
		return
	}

	err = cfg.authorizeMember(shelfLocation.ID, *r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to gets movies at the location of that shelf", err)
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
			Format:      dbMovie.Format,
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

	// Validate user is authorized to get movies at the location of requested movie.
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
		Format:      dbMovie.Format,
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
		Format:      dbMovie.Format,
		ReleaseDate: dbMovie.ReleaseDate,
		CreatedAt:   dbMovie.CreatedAt,
		UpdatedAt:   dbMovie.UpdatedAt,
		ShelfID:     dbMovie.ShelfID,
	})
}

func (cfg *apiConfig) handlerMoviesGetByLocation(w http.ResponseWriter, r *http.Request) {
	locationIDString := r.PathValue("location_id")
	if locationIDString == "" {
		respondWithError(w, http.StatusBadRequest, "No location id was provided", fmt.Errorf("no location id was provided"))
		return
	}

	locationID, err := uuid.Parse(locationIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid location ID", err)
		return
	}

	// Validate user is permitted to get movies for the location.
	err = cfg.authorizeMember(locationID, *r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to gets movies for this location", err)
		return
	}

	dbMovies, err := cfg.db.GetMoviesByLocation(r.Context(), locationID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No movies found for that location", err)
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
			Format:      dbMovie.Format,
			ReleaseDate: dbMovie.ReleaseDate,
			CreatedAt:   dbMovie.CreatedAt,
			UpdatedAt:   dbMovie.UpdatedAt,
			ShelfID:     dbMovie.ShelfID,
		})
	}

	respondWithJSON(w, http.StatusOK, movies)
}

func (cfg *apiConfig) handlerSearchMovies(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		LocationID string `json:"location_id"`
		Query      string `json:"query"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	locationIDString := requestBody.LocationID
	if locationIDString == "" {
		respondWithError(w, http.StatusBadRequest, "No location id was provided", fmt.Errorf("no location id was provided"))
		return
	}

	locationID, err := uuid.Parse(locationIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid location ID", err)
		return
	}

	// Validate user is permitted to search movies for the location.
	err = cfg.authorizeMember(locationID, *r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to search movies for this location", err)
		return
	}

	query := requestBody.Query
	if query == "" {
		respondWithError(w, http.StatusBadRequest, "No search query was provided", fmt.Errorf("no search query was provided"))
		return
	}

	dbMovies, err := cfg.db.SearchMovies(r.Context(), database.SearchMoviesParams{
		WebsearchToTsquery: query,
		ID:                 locationID,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No movies found for that location", err)
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
			Format:      dbMovie.Format,
			ReleaseDate: dbMovie.ReleaseDate,
			CreatedAt:   dbMovie.CreatedAt,
			UpdatedAt:   dbMovie.UpdatedAt,
			ShelfID:     dbMovie.ShelfID,
		})
	}

	respondWithJSON(w, http.StatusOK, movies)
}
