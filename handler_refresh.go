package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Rodabaugh/digitalshelf/internal/auth"
)

func (apiCfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	fmt.Println("refreshToken: ", refreshToken)
	user, err := apiCfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	fmt.Printf("user: %v\n", user)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user for refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		apiCfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}

func (apiCfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	_, err = apiCfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke session", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (apiCfg *apiConfig) handlerRevokeSessions(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization header is required", nil)
		return
	}

	tokenStr := authHeader[len("Bearer "):]
	userID, err := auth.ValidateJWT(tokenStr, apiCfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	_, err = apiCfg.db.RevokeRefreshTokenByUserId(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke refresh tokens", err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Sessions revoked successfully"})
}
