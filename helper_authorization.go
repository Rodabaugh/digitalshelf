package main

import (
	"fmt"
	"net/http"

	"github.com/Rodabaugh/digitalshelf/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) authorizeOwner(locationID uuid.UUID, r http.Request) (err error) {
	if len(locationID) == 0 {
		return fmt.Errorf("location ID is required")
	}

	userID, err := cfg.getRequesterID(r)
	if err != nil {
		return fmt.Errorf("unable to get requester ID: %w", err)
	}

	location, err := cfg.db.GetLocationByID(r.Context(), locationID)
	if err != nil {
		return fmt.Errorf("unable to get location: %w", err)
	}

	if location.OwnerID != userID {
		return fmt.Errorf("user is not the owner of this location")
	}

	return nil
}

func (cfg *apiConfig) authorizeUser(locationID uuid.UUID, r http.Request) (err error) {
	if len(locationID) == 0 {
		return fmt.Errorf("location ID is required")
	}

	userID, err := cfg.getRequesterID(r)
	if err != nil {
		return fmt.Errorf("unable to get requester ID: %w", err)
	}

	locationMembers, err := cfg.db.GetLocationMembers(r.Context(), locationID)
	if err != nil {
		return fmt.Errorf("unable to get location members: %w", err)
	}
	for _, member := range locationMembers {
		if member.ID == userID {
			return nil
		}
	}
	return fmt.Errorf("user is not a member of this location")
}

func (cfg *apiConfig) authorizeInvited(locationID uuid.UUID, r http.Request) (err error) {
	if len(locationID) == 0 {
		return fmt.Errorf("location ID is required")
	}

	userID, err := cfg.getRequesterID(r)
	if err != nil {
		return fmt.Errorf("unable to get requester ID: %w", err)
	}

	locationInvites, err := cfg.db.GetLocationInvites(r.Context(), locationID)
	if err != nil {
		return fmt.Errorf("unable to get location invites: %w", err)
	}

	for _, invite := range locationInvites {
		if invite.ID == userID {
			return nil
		}
	}
	return fmt.Errorf("user is not invited to this location")
}

func (cfg *apiConfig) getRequesterID(r http.Request) (uuid.UUID, error) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return uuid.Nil, fmt.Errorf("unable to get bearer token: %w", err)
	}
	userID, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		return uuid.Nil, fmt.Errorf("unable to validate JWT: %w", err)
	}
	return userID, nil
}
