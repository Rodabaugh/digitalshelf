package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLocationsGet(w http.ResponseWriter, r *http.Request) {
	dbLocations, err := cfg.db.GetLocations(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get locations from database", err)
		return
	}

	locations := []Location{}

	for _, dbLocation := range dbLocations {
		locations = append(locations, Location{
			ID:        dbLocation.ID,
			Name:      dbLocation.Name,
			OwnerID:   dbLocation.OwnerID,
			CreatedAt: dbLocation.CreatedAt,
			UpdatedAt: dbLocation.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, locations)
}

func (cfg *apiConfig) handlerLocationsGetByOwner(w http.ResponseWriter, r *http.Request) {
	ownerIDString := r.URL.Query().Get("owner_id")
	if ownerIDString == "" {
		respondWithError(w, http.StatusBadRequest, "No owner_id was provided", fmt.Errorf("no owner_id was provided"))
		return
	}

	ownerID, err := uuid.Parse(ownerIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	dbLocations, err := cfg.db.GetLocationsByOwner(r.Context(), ownerID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No locations owned by that user", err)
		return
	}

	locations := []Location{}

	for _, dbLocation := range dbLocations {
		locations = append(locations, Location{
			ID:        dbLocation.ID,
			Name:      dbLocation.Name,
			OwnerID:   dbLocation.OwnerID,
			CreatedAt: dbLocation.CreatedAt,
			UpdatedAt: dbLocation.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, locations)
}

func (cfg *apiConfig) handlerLocationsGetByID(w http.ResponseWriter, r *http.Request) {
	locationIDString := r.PathValue("location_id")
	if locationIDString == "" {
		respondWithError(w, http.StatusBadRequest, "No location id was provided", fmt.Errorf("no location id was provided"))
		return
	}

	locationID, err := uuid.Parse(locationIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	dbLocation, err := cfg.db.GetLocationByID(r.Context(), locationID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Location not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Location{
		ID:        dbLocation.ID,
		Name:      dbLocation.Name,
		OwnerID:   dbLocation.OwnerID,
		CreatedAt: dbLocation.CreatedAt,
		UpdatedAt: dbLocation.UpdatedAt,
	})
}
