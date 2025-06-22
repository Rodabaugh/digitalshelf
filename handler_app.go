package main

import (
	"net/http"
	"fmt"

	"github.com/Rodabaugh/digitalshelf/internal/auth"
	"github.com/google/uuid"
)

type AppState struct {
	userID	string
}

func (cfg *apiConfig) webApp(w http.ResponseWriter, r *http.Request) {
	var cookieUserID uuid.UUID

	accessTokenCookie, err := r.Cookie("accessToken")

	if err != nil {
		if err != http.ErrNoCookie{
			fmt.Printf("Error reading cookie: %v\n", err)
		}
	} else {
		cookieUserID, err = auth.ValidateJWT(accessTokenCookie.Value, cfg.jwtSecret)
		if err != nil{
			fmt.Printf("Error validating JWT: %v\n", err)
		}
	}

	appState := AppState{
		userID: cookieUserID.String(),
	}

	MainPage(&appState).Render(r.Context(), w)
}
