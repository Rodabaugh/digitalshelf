package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsersGet(w http.ResponseWriter, r *http.Request) {
	dbUsers, err := cfg.db.GetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get users from database", err)
		return
	}

	users := []User{}

	for _, dbUser := range dbUsers {
		users = append(users, User{
			ID:        dbUser.ID,
			Name:      dbUser.Name,
			Email:     dbUser.Email,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, users)
}

func (cfg *apiConfig) handlerUsersGetByEmail(w http.ResponseWriter, r *http.Request) {
	userEmail := r.URL.Query().Get("email")
	if userEmail == "" {
		respondWithError(w, http.StatusBadRequest, "No email was provided", fmt.Errorf("no email was provided"))
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(r.Context(), userEmail)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	})
}

func (cfg *apiConfig) handlerUserGetByID(w http.ResponseWriter, r *http.Request) {
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

	dbUser, err := cfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	})
}
