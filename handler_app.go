package main

import (
	"net/http"
	"encoding/json"
	"fmt"

	"github.com/Rodabaugh/digitalshelf/internal/auth"
	"github.com/Rodabaugh/digitalshelf/internal/database"
	"github.com/google/uuid"
)

type AppState struct {
	userID	string
}

func (cfg *apiConfig) webApp(w http.ResponseWriter, r *http.Request) {
	cookieUserID := cfg.getRequestUserID(r)

	appState := AppState{
		userID: cookieUserID.String(),
	}

	MainPage(&appState).Render(r.Context(), w)
}


func (cfg *apiConfig) appGetUserLocations(w http.ResponseWriter, r *http.Request) {
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

	requestUserID := cfg.getRequestUserID(r)

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

	Locations(userLocations).Render(r.Context(), w)
}

func (cfg *apiConfig) appCreateLocation (w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name    string    `json:"name"`
	}

	type response struct {
		Location
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Was unable to decode parameters", err)
		return
	}

	requestUserID := cfg.getRequestUserID(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get request user ID", err)
		return
	}

	dbUser, err := cfg.db.GetUserByID(r.Context(), requestUserID)
	if err != nil{
		fmt.Printf("Unable to get user: %v\n", err)
		return
	}

	location, err := cfg.db.CreateLocation(r.Context(), database.CreateLocationParams{
		Name:    params.Name,
		OwnerID: requestUserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create location", err)
		return
	}

	newLocationMemberParams := database.AddLocationMemberParams{	
		LocationID: location.ID,
		UserID: dbUser.ID,
	}

	_, err = cfg.db.AddLocationMember(r.Context(), newLocationMemberParams)
	if err != nil {
		fmt.Printf("Failed to add user to location: %v\n", err)
	}

	AppSuccessReply().Render(r.Context(), w)
}

func (cfg apiConfig) getRequestUserID(r *http.Request) (uuid.UUID){
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

	return cookieUserID
}
