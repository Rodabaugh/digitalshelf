package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Rodabaugh/digitalshelf/internal/database"
	"github.com/google/uuid"
)

type NewLocationInvite struct {
	LocationID uuid.UUID `json:"location_id"`
	UserID     uuid.UUID `json:"user_id"`
	InvitedAt  time.Time `json:"invited_at"`
}

type LocationInvite struct {
	LocationID uuid.UUID `json:"location_id"`
	UserID     uuid.UUID `json:"userID"`
	UserName   string    `json:"user_name"`
	UserEmail  string    `json:"user_email"`
	InvitedAt  time.Time `json:"invited_at"`
}

type UserInvite struct {
	UserID       uuid.UUID `json:"userID"`
	LocationID   uuid.UUID `json:"location_id"`
	LocationName string    `json:"location_name"`
	OwnerID      uuid.UUID `json:"owner_id"`
	InvitedAt    time.Time `json:"invited_at"`
}

func (cfg *apiConfig) handlerAddLocationInvite(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserID string `json:"user_id"`
	}

	type response struct {
		NewLocationInvite
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Was unable to decode parameters", err)
		return
	}

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

	// Validate that the user is permitted to add invites to this location.
	err = cfg.authorizeOwner(locationID, *r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to add invites to this location", err)
		return
	}

	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	LocationInvite, err := cfg.db.AddLocationInvite(r.Context(), database.AddLocationInviteParams{
		LocationID: locationID,
		UserID:     userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to add user to location", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		NewLocationInvite: NewLocationInvite{
			LocationID: LocationInvite.LocationID,
			UserID:     LocationInvite.UserID,
			InvitedAt:  time.Now(),
		},
	})
}

func (cfg *apiConfig) handlerGetLocationInvites(w http.ResponseWriter, r *http.Request) {
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

	// Validate that the user is permitted to get the invites for this location.
	err = cfg.authorizeOwner(locationID, *r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to view invites for this location", err)
		return
	}

	dbLocationInvites, err := cfg.db.GetLocationInvites(r.Context(), locationID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get invites for location", err)
		return
	}

	locationInvites := []LocationInvite{}

	for _, locationInvite := range dbLocationInvites {
		locationInvites = append(locationInvites, LocationInvite{
			LocationID: locationInvite.LocationID,
			UserID:     locationInvite.ID,
			UserName:   locationInvite.Name,
			UserEmail:  locationInvite.Email,
			InvitedAt:  locationInvite.InvitedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, locationInvites)
}

func (cfg *apiConfig) handlerGetUserInvites(w http.ResponseWriter, r *http.Request) {
	userIDString := r.PathValue("user_id")
	if userIDString == "" {
		respondWithError(w, http.StatusBadRequest, "No user id was provided", fmt.Errorf("no user id was provided"))
		return
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	// Validate that the user is permitted to get the invites for this user.
	requesterID, err := cfg.getRequesterID(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to view invites for this user", err)
		return
	}

	if userID != requesterID {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to view invites for this user", err)
	}

	dbUserInvites, err := cfg.db.GetUserInvites(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get invites for user", err)
		return
	}

	userInvites := []UserInvite{}

	for _, userInvite := range dbUserInvites {
		userInvites = append(userInvites, UserInvite{
			UserID:       userInvite.UserID,
			LocationID:   userInvite.ID,
			LocationName: userInvite.Name,
			OwnerID:      userInvite.OwnerID,
			InvitedAt:    userInvite.InvitedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, userInvites)
}

func (cfg *apiConfig) handlerRemoveLocationInvite(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserID string `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Was unable to decode parameters", err)
		return
	}

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

	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	// Check if the requester is the owner or an invited user.
	// The below calls returns a nil error if the user is authorized.
	isOwner := cfg.authorizeOwner(locationID, *r)
	isInvited := cfg.authorizeInvited(locationID, *r)
	if isOwner != nil || isInvited != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to remove this invite", fmt.Errorf("user is not authorized to remove this invite"))
		return
	}

	err = cfg.db.RemoveLocationInvite(r.Context(), database.RemoveLocationInviteParams{
		UserID:     userID,
		LocationID: locationID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete location invite", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
