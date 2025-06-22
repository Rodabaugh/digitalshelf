package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Rodabaugh/digitalshelf/internal/auth"
	"github.com/Rodabaugh/digitalshelf/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Was unable to decode parameters", err)
		return
	}

	if len(params.Email) == 0 {
		respondWithError(w, http.StatusBadRequest, "An email address must be provided", err)
	}

	if len(params.Password) == 0 {
		respondWithError(w, http.StatusBadRequest, "A password is required", err)
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		log.Printf("Failed to hash password %s", err)
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Name:           params.Name,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create user", err)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		respondWithJSON(w, http.StatusCreated, response{
			User: User{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
			},
		})
	} else {
		RegisterSuccess().Render(r.Context(), w)
	}

	fmt.Printf("User %s Created\n", user.Email)
}
