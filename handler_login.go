package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Rodabaugh/digitalshelf/internal/auth"
	"github.com/Rodabaugh/digitalshelf/internal/database"
)

func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	parms := parameters{}

	w.Header().Set("Content-Type", "application/json")

	err := decoder.Decode(&parms)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		log.Printf("Error decoding parameters: %s", err)
	}

	if len(parms.Email) == 0 {
		respondWithError(w, http.StatusBadRequest, "An email address must be provided", err)
	}

	if len(parms.Password) == 0 {
		respondWithError(w, http.StatusBadRequest, "A password is required", err)
	}

	user, err := apiCfg.db.GetUserByEmail(r.Context(), parms.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Was unable to find user.", err)
		log.Printf("Unable to get user from database: %s", err)
	}

	err = auth.CheckPasswordHash(parms.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		apiCfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	_, err = apiCfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Name:      user.Name,
			Email:     user.Email,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
	fmt.Println("Login successful")
}
