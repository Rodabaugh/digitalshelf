package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rodabaugh/digitalshelf/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerShowsGet(w http.ResponseWriter, r *http.Request) {
	dbShows, err := cfg.db.GetShows(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get shows from database", err)
		return
	}

	shows := []Show{}

	for _, dbShow := range dbShows {
		shows = append(shows, Show{
			ID:          dbShow.ID,
			Title:       dbShow.Title,
			Genre:       dbShow.Genre,
			Actors:      dbShow.Actors,
			Writer:      dbShow.Writer,
			Director:    dbShow.Director,
			Barcode:     dbShow.Barcode,
			ReleaseDate: dbShow.ReleaseDate,
			CreatedAt:   dbShow.CreatedAt,
			UpdatedAt:   dbShow.UpdatedAt,
			ShelfID:     dbShow.ShelfID,
		})
	}

	respondWithJSON(w, http.StatusOK, shows)
}

func (cfg *apiConfig) handlerShowsGetByShelf(w http.ResponseWriter, r *http.Request) {
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

	// Validate user is authorized to get shows at the location of shelf.
	shelfLocation, err := cfg.db.GetShelfLocation(r.Context(), shelfID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get shelf location", err)
		return
	}

	err = cfg.authorizeMember(shelfLocation.ID, *r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to get shows at the location of that shelf", err)
		return
	}

	dbShows, err := cfg.db.GetShowsByShelf(r.Context(), shelfID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No shows found for that shelf", err)
		return
	}

	shows := []Show{}

	for _, dbShow := range dbShows {
		shows = append(shows, Show{
			ID:          dbShow.ID,
			Title:       dbShow.Title,
			Season:      dbShow.Season,
			Genre:       dbShow.Genre,
			Actors:      dbShow.Actors,
			Writer:      dbShow.Writer,
			Director:    dbShow.Director,
			Barcode:     dbShow.Barcode,
			ReleaseDate: dbShow.ReleaseDate,
			CreatedAt:   dbShow.CreatedAt,
			UpdatedAt:   dbShow.UpdatedAt,
			ShelfID:     dbShow.ShelfID,
		})
	}

	respondWithJSON(w, http.StatusOK, shows)
}

func (cfg *apiConfig) handlerShowGetByID(w http.ResponseWriter, r *http.Request) {
	showIDString := r.PathValue("show_id")
	if showIDString == "" {
		respondWithError(w, http.StatusBadRequest, "No show id was provided", fmt.Errorf("no show id was provided"))
		return
	}

	showID, err := uuid.Parse(showIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid show ID", err)
		return
	}

	// Validate user is authorized to get shows at the location of requested show.
	showLocation, err := cfg.db.GetShowLocation(r.Context(), showID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get show location", err)
		return
	}

	err = cfg.authorizeMember(showLocation.ID, *r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to get shows at this location", err)
		return
	}

	dbShow, err := cfg.db.GetShowByID(r.Context(), showID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Show not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Show{
		ID:          dbShow.ID,
		Title:       dbShow.Title,
		Season:      dbShow.Season,
		Genre:       dbShow.Genre,
		Actors:      dbShow.Actors,
		Writer:      dbShow.Writer,
		Director:    dbShow.Director,
		Barcode:     dbShow.Barcode,
		ReleaseDate: dbShow.ReleaseDate,
		CreatedAt:   dbShow.CreatedAt,
		UpdatedAt:   dbShow.UpdatedAt,
		ShelfID:     dbShow.ShelfID,
	})
}

func (cfg *apiConfig) handlerGetShowByBarcode(w http.ResponseWriter, r *http.Request) {
	barcode := r.PathValue("barcode")
	fmt.Println(barcode)
	if barcode == "" {
		respondWithError(w, http.StatusBadRequest, "No barcode was provided", fmt.Errorf("no barcode was provided"))
		return
	}

	dbShow, err := cfg.db.GetShowByBarcode(r.Context(), barcode)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Show not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Show{
		ID:          dbShow.ID,
		Title:       dbShow.Title,
		Season:      dbShow.Season,
		Genre:       dbShow.Genre,
		Actors:      dbShow.Actors,
		Writer:      dbShow.Writer,
		Director:    dbShow.Director,
		Barcode:     dbShow.Barcode,
		ReleaseDate: dbShow.ReleaseDate,
		CreatedAt:   dbShow.CreatedAt,
		UpdatedAt:   dbShow.UpdatedAt,
		ShelfID:     dbShow.ShelfID,
	})
}

func (cfg *apiConfig) handlerShowsGetByLocation(w http.ResponseWriter, r *http.Request) {
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

	// Validate user is permitted to get shows for the location.
	err = cfg.authorizeMember(locationID, *r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to get shows for this location", err)
		return
	}

	dbShows, err := cfg.db.GetShowsByLocation(r.Context(), locationID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No shows found for that location", err)
		return
	}

	shows := []Show{}

	for _, dbShow := range dbShows {
		shows = append(shows, Show{
			ID:          dbShow.ID,
			Title:       dbShow.Title,
			Season:      dbShow.Season,
			Genre:       dbShow.Genre,
			Actors:      dbShow.Actors,
			Writer:      dbShow.Writer,
			Director:    dbShow.Director,
			Barcode:     dbShow.Barcode,
			ReleaseDate: dbShow.ReleaseDate,
			CreatedAt:   dbShow.CreatedAt,
			UpdatedAt:   dbShow.UpdatedAt,
			ShelfID:     dbShow.ShelfID,
		})
	}

	respondWithJSON(w, http.StatusOK, shows)
}

func (cfg *apiConfig) handlerSearchShows(w http.ResponseWriter, r *http.Request) {
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

	// Validate user is permitted to search shows for the location.
	err = cfg.authorizeMember(locationID, *r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to search shows for this location", err)
		return
	}

	query := requestBody.Query
	if query == "" {
		respondWithError(w, http.StatusBadRequest, "No search query was provided", fmt.Errorf("no search query was provided"))
		return
	}

	dbShows, err := cfg.db.SearchShows(r.Context(), database.SearchShowsParams{
		WebsearchToTsquery: query,
		ID:                 locationID,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No shows found for that location", err)
		return
	}

	shows := []Show{}

	for _, dbShow := range dbShows {
		shows = append(shows, Show{
			ID:          dbShow.ID,
			Title:       dbShow.Title,
			Season:      dbShow.Season,
			Genre:       dbShow.Genre,
			Actors:      dbShow.Actors,
			Writer:      dbShow.Writer,
			Director:    dbShow.Director,
			Barcode:     dbShow.Barcode,
			ReleaseDate: dbShow.ReleaseDate,
			CreatedAt:   dbShow.CreatedAt,
			UpdatedAt:   dbShow.UpdatedAt,
			ShelfID:     dbShow.ShelfID,
		})
	}

	respondWithJSON(w, http.StatusOK, shows)
}
