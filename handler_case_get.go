package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCaseGet(w http.ResponseWriter, r *http.Request) {
	dbCases, err := cfg.db.GetCases(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get cases from database", err)
		return
	}

	itemCases := []Case{}

	for _, dbCase := range dbCases {
		itemCases = append(itemCases, Case{
			ID:         dbCase.ID,
			Name:       dbCase.Name,
			LocationID: dbCase.LocationID,
			CreatedAt:  dbCase.CreatedAt,
			UpdatedAt:  dbCase.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, itemCases)
}

func (cfg *apiConfig) handlerCasesGetByLocation(w http.ResponseWriter, r *http.Request) {
	locationIDString := r.PathValue("location_id")
	if locationIDString == "" {
		respondWithError(w, http.StatusBadRequest, "No location was provided", fmt.Errorf("no location_id was provided"))
		return
	}

	locationID, err := uuid.Parse(locationIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid location ID", err)
		return
	}

	dbCases, err := cfg.db.GetCasesByLocation(r.Context(), locationID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No cases found for that location", err)
		return
	}

	itemCases := []Case{}

	for _, dbCase := range dbCases {
		itemCases = append(itemCases, Case{
			ID:         dbCase.ID,
			Name:       dbCase.Name,
			LocationID: dbCase.LocationID,
			CreatedAt:  dbCase.CreatedAt,
			UpdatedAt:  dbCase.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, itemCases)
}

func (cfg *apiConfig) handlerCaseGetByID(w http.ResponseWriter, r *http.Request) {
	caseIDString := r.PathValue("case_id")
	if caseIDString == "" {
		respondWithError(w, http.StatusBadRequest, "No case id was provided", fmt.Errorf("no case id was provided"))
		return
	}

	caseID, err := uuid.Parse(caseIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid case ID", err)
		return
	}

	dbCase, err := cfg.db.GetCaseByID(r.Context(), caseID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Case not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Case{
		ID:         dbCase.ID,
		Name:       dbCase.Name,
		LocationID: dbCase.LocationID,
		CreatedAt:  dbCase.CreatedAt,
		UpdatedAt:  dbCase.UpdatedAt,
	})
}
