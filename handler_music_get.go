package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rodabaugh/digitalshelf/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerMusicGet(w http.ResponseWriter, r *http.Request) {
	dbMusic, err := cfg.db.GetMusic(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get music from database", err)
		return
	}

	music := []Music{}

	for _, dbM := range dbMusic {
		music = append(music, Music{
			ID:          dbM.ID,
			Title:       dbM.Title,
			Artist:      dbM.Artist,
			Genre:       dbM.Genre,
			Barcode:     dbM.Barcode,
			Format:      dbM.Format,
			ShelfID:     dbM.ShelfID,
			ReleaseDate: dbM.ReleaseDate,
			CreatedAt:   dbM.CreatedAt,
			UpdatedAt:   dbM.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, music)
}

func (cfg *apiConfig) handlerMusicGetByShelf(w http.ResponseWriter, r *http.Request) {
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

	// Validate user is authorized to get music at the location of shelf.
	shelfLocation, err := cfg.db.GetShelfLocation(r.Context(), shelfID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get shelf location", err)
		return
	}

	err = cfg.authorizeMember(shelfLocation.ID, *r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to get music at the location of that shelf", err)
		return
	}

	dbMusic, err := cfg.db.GetMusicByShelf(r.Context(), shelfID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No music found for that shelf", err)
		return
	}

	music := []Music{}

	for _, dbM := range dbMusic {
		music = append(music, Music{
			ID:          dbM.ID,
			Title:       dbM.Title,
			Artist:      dbM.Artist,
			Genre:       dbM.Genre,
			Barcode:     dbM.Barcode,
			Format:      dbM.Format,
			ShelfID:     dbM.ShelfID,
			ReleaseDate: dbM.ReleaseDate,
			CreatedAt:   dbM.CreatedAt,
			UpdatedAt:   dbM.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, music)
}

func (cfg *apiConfig) handlerMusicGetByID(w http.ResponseWriter, r *http.Request) {
	musicIDString := r.PathValue("music_id")
	if musicIDString == "" {
		respondWithError(w, http.StatusBadRequest, "No music id was provided", fmt.Errorf("no music id was provided"))
		return
	}

	musicID, err := uuid.Parse(musicIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid music ID", err)
		return
	}

	// Validate user is authorized to get music at the location of requested music.
	musicLocation, err := cfg.db.GetMusicLocation(r.Context(), musicID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get music location", err)
		return
	}

	err = cfg.authorizeMember(musicLocation.ID, *r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to get music at this location", err)
		return
	}

	dbMusic, err := cfg.db.GetMusicByID(r.Context(), musicID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Music not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Music{
		ID:          dbMusic.ID,
		Title:       dbMusic.Title,
		Artist:      dbMusic.Artist,
		Genre:       dbMusic.Genre,
		Barcode:     dbMusic.Barcode,
		Format:      dbMusic.Format,
		ShelfID:     dbMusic.ShelfID,
		ReleaseDate: dbMusic.ReleaseDate,
		CreatedAt:   dbMusic.CreatedAt,
		UpdatedAt:   dbMusic.UpdatedAt,
	})
}

func (cfg *apiConfig) handlerGetMusicByBarcode(w http.ResponseWriter, r *http.Request) {
	barcode := r.PathValue("barcode")
	fmt.Println(barcode)
	if barcode == "" {
		respondWithError(w, http.StatusBadRequest, "No barcode was provided", fmt.Errorf("no barcode was provided"))
		return
	}

	dbMusic, err := cfg.db.GetMusicByBarcode(r.Context(), barcode)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Music not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Music{
		ID:          dbMusic.ID,
		Title:       dbMusic.Title,
		Artist:      dbMusic.Artist,
		Genre:       dbMusic.Genre,
		Barcode:     dbMusic.Barcode,
		Format:      dbMusic.Format,
		ShelfID:     dbMusic.ShelfID,
		ReleaseDate: dbMusic.ReleaseDate,
		CreatedAt:   dbMusic.CreatedAt,
		UpdatedAt:   dbMusic.UpdatedAt,
	})
}

func (cfg *apiConfig) handlerMusicGetByLocation(w http.ResponseWriter, r *http.Request) {
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

	// Validate user is permitted to get music for the location.
	err = cfg.authorizeMember(locationID, *r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to get music for this location", err)
		return
	}

	dbMusic, err := cfg.db.GetMusicByLocation(r.Context(), locationID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No music found for that location", err)
		return
	}

	music := []Music{}

	for _, dbM := range dbMusic {
		music = append(music, Music{
			ID:          dbM.ID,
			Title:       dbM.Title,
			Artist:      dbM.Artist,
			Genre:       dbM.Genre,
			Barcode:     dbM.Barcode,
			Format:      dbM.Format,
			ShelfID:     dbM.ShelfID,
			ReleaseDate: dbM.ReleaseDate,
			CreatedAt:   dbM.CreatedAt,
			UpdatedAt:   dbM.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, music)
}

func (cfg *apiConfig) handlerSearchMusic(w http.ResponseWriter, r *http.Request) {
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

	// Validate user is permitted to search music for the location.
	err = cfg.authorizeMember(locationID, *r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to search music for this location", err)
		return
	}

	query := requestBody.Query
	if query == "" {
		respondWithError(w, http.StatusBadRequest, "No search query was provided", fmt.Errorf("no search query was provided"))
		return
	}

	dbMusic, err := cfg.db.SearchMusic(r.Context(), database.SearchMusicParams{
		WebsearchToTsquery: query,
		ID:                 locationID,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No music found for that location", err)
		return
	}

	music := []Music{}

	for _, dbM := range dbMusic {
		music = append(music, Music{
			ID:          dbM.ID,
			Title:       dbM.Title,
			Artist:      dbM.Artist,
			Genre:       dbM.Genre,
			Barcode:     dbM.Barcode,
			Format:      dbM.Format,
			ShelfID:     dbM.ShelfID,
			ReleaseDate: dbM.ReleaseDate,
			CreatedAt:   dbM.CreatedAt,
			UpdatedAt:   dbM.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, music)
}
