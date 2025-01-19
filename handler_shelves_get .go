package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerShelvesGet(w http.ResponseWriter, r *http.Request) {
	dbShelves, err := cfg.db.GetShelves(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get shelves from database", err)
		return
	}

	shelves := []Shelf{}

	for _, dbShelf := range dbShelves {
		shelves = append(shelves, Shelf{
			ID:        dbShelf.ID,
			Name:      dbShelf.Name,
			CaseID:    dbShelf.CaseID,
			CreatedAt: dbShelf.CreatedAt,
			UpdatedAt: dbShelf.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, shelves)
}

func (cfg *apiConfig) handlerShelvesGetByCase(w http.ResponseWriter, r *http.Request) {
	caseIDString := r.PathValue("case_id")
	if caseIDString == "" {
		respondWithError(w, http.StatusBadRequest, "No case id was provided", fmt.Errorf("no case_id was provided"))
		return
	}

	caseID, err := uuid.Parse(caseIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid case ID", err)
		return
	}

	// Validate user is authorized to get shelves at the case's location.
	caseLocation, err := cfg.db.GetCaseLocation(r.Context(), caseID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get case location", err)
		return
	}

	if err := cfg.authorizeMember(caseLocation.ID, *r); err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to get shelves at this case", err)
		return
	}

	dbShelves, err := cfg.db.GetShelvesByCase(r.Context(), caseID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No shelves found for that case", err)
		return
	}

	shelves := []Shelf{}

	for _, dbShelf := range dbShelves {
		shelves = append(shelves, Shelf{
			ID:        dbShelf.ID,
			Name:      dbShelf.Name,
			CaseID:    dbShelf.CaseID,
			CreatedAt: dbShelf.CreatedAt,
			UpdatedAt: dbShelf.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, shelves)
}

func (cfg *apiConfig) handlerShelfGetByID(w http.ResponseWriter, r *http.Request) {
	shelfIDString := r.PathValue("shelf_id")
	if shelfIDString == "" {
		respondWithError(w, http.StatusBadRequest, "No shelf id was provided", fmt.Errorf("no shelf id was provided"))
		return
	}

	shelfID, err := uuid.Parse(shelfIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid shelf ID", err)
		return
	}

	// Validate user is authorized to get shelves at the shelf's location.
	shelfLocation, err := cfg.db.GetShelfLocation(r.Context(), shelfID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get shelf location", err)
		return
	}

	if err := cfg.authorizeMember(shelfLocation.ID, *r); err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to get this shelf", err)
		return
	}

	dbShelf, err := cfg.db.GetShelfByID(r.Context(), shelfID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Shelf not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Shelf{
		ID:        dbShelf.ID,
		Name:      dbShelf.Name,
		CaseID:    dbShelf.CaseID,
		CreatedAt: dbShelf.CreatedAt,
		UpdatedAt: dbShelf.UpdatedAt,
	})
}
