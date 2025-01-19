package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Rodabaugh/digitalshelf/internal/database"
	"github.com/google/uuid"
)

type NewLocationUser struct {
	LocationID uuid.UUID `json:"location_id"`
	UserID     uuid.UUID `json:"user_id"`
	JoinedAt   time.Time `json:"joined_at"`
}

type LocationMember struct {
	LocationID uuid.UUID `json:"location_id"`
	UserID     uuid.UUID `json:"userID"`
	UserName   string    `json:"user_name"`
	UserEmail  string    `json:"user_email"`
	JoinedAt   time.Time `json:"joined_at"`
}

type UserLocation struct {
	UserID       uuid.UUID `json:"userID"`
	LocationID   uuid.UUID `json:"location_id"`
	LocationName string    `json:"location_name"`
	OwnerID      uuid.UUID `json:"owner_id"`
	JoinedAt     time.Time `json:"joined_at"`
}

func (cfg *apiConfig) handlerAddLocationMember(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserID string `json:"user_id"`
	}

	type response struct {
		NewLocationUser
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
	if isOwner != nil && isInvited != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to join location", fmt.Errorf("user is not authorized to join location"))
		return
	}

	locationUser, err := cfg.db.AddLocationMember(r.Context(), database.AddLocationMemberParams{
		LocationID: locationID,
		UserID:     userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to add user to location", err)
		return
	}

	if isInvited == nil {
		// If the user was invited, we need to remove the invite.
		err = cfg.db.RemoveLocationInvite(r.Context(), database.RemoveLocationInviteParams{
			UserID:     userID,
			LocationID: locationID,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to remove location invite: %v", err)
		}
	}

	respondWithJSON(w, http.StatusCreated, response{
		NewLocationUser: NewLocationUser{
			LocationID: locationUser.LocationID,
			UserID:     locationUser.UserID,
			JoinedAt:   locationUser.JoinedAt,
		},
	})
}

func (cfg *apiConfig) handlerGetLocationMembers(w http.ResponseWriter, r *http.Request) {
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

	// Validate user is the owner of the location.
	if cfg.authorizeOwner(locationID, *r) != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to view members of this location", err)
		return
	}

	dbLocationMembers, err := cfg.db.GetLocationMembers(r.Context(), locationID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get members for location", err)
		return
	}

	locationMembers := []LocationMember{}

	for _, locationMember := range dbLocationMembers {
		locationMembers = append(locationMembers, LocationMember{
			LocationID: locationMember.LocationID,
			UserID:     locationMember.ID,
			UserName:   locationMember.Name,
			UserEmail:  locationMember.Email,
			JoinedAt:   locationMember.JoinedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, locationMembers)
}

func (cfg *apiConfig) handlerGetUserLocations(w http.ResponseWriter, r *http.Request) {
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

	// Validate the user is authorized to view locations for this user.
	requestUserID, err := cfg.getRequesterID(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get requester ID", err)
		return
	}
	if userID != requestUserID {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to view locations for this user", err)
		return
	}

	dbUserLocations, err := cfg.db.GetUserLocations(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get locations for user", err)
		return
	}

	userLocations := []UserLocation{}

	for _, userLocation := range dbUserLocations {
		userLocations = append(userLocations, UserLocation{
			UserID:       userLocation.UserID,
			LocationID:   userLocation.ID,
			LocationName: userLocation.Name,
			OwnerID:      userLocation.OwnerID,
			JoinedAt:     userLocation.JoinedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, userLocations)
}

func (cfg *apiConfig) handlerRemoveLocationMember(w http.ResponseWriter, r *http.Request) {
	locationIDString := r.PathValue("location_id")
	if locationIDString == "" {
		respondWithError(w, http.StatusBadRequest, "No location id was provided", fmt.Errorf("no location id was provided"))
		return
	}

	userIDString := r.PathValue("user_id")
	if userIDString == "" {
		respondWithError(w, http.StatusBadRequest, "No user id was provided", fmt.Errorf("no user id was provided"))
	}
	println(userIDString)

	locationID, err := uuid.Parse(locationIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid location ID", err)
		return
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	// Check if the requester is the owner or an invited user.
	// The below calls returns a nil error if the user is authorized.
	isOwner := cfg.authorizeOwner(locationID, *r)
	isMember := cfg.authorizeMember(locationID, *r)
	if isOwner != nil && isMember != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to remove this user", fmt.Errorf("user is not authorized to remove this user"))
		return
	}

	err = cfg.db.RemoveLocationMember(r.Context(), database.RemoveLocationMemberParams{
		UserID:     userID,
		LocationID: locationID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete user location membership", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
