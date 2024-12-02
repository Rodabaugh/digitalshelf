package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Rodabaugh/digitalshelf/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	platform string
	db       *database.Queries
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform != "dev" && platform != "prod" {
		log.Fatal("PLATFORM must be set to either dev or prod")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		platform: platform,
		db:       dbQueries,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", readinessEndpoint)

	mux.HandleFunc("POST /users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("POST /locations", apiCfg.handlerLocationsCreate)
	mux.HandleFunc("POST /cases", apiCfg.handlerCasesCreate)
	mux.HandleFunc("GET /cases", apiCfg.handlerCaseGet)

	mux.HandleFunc("GET /users", apiCfg.handlerUsersGet)
	mux.HandleFunc("GET /users/{user_id}", apiCfg.handlerUserGetByID)
	mux.HandleFunc("GET /users/{user_id}/locations", apiCfg.handlerGetUserLocations)
	mux.HandleFunc("GET /locations", apiCfg.handlerLocationsGet)
	mux.HandleFunc("GET /locations/{location_id}", apiCfg.handlerLocationsGetByID)
	mux.HandleFunc("GET /locations/{location_id}/members", apiCfg.handlerGetLocationMembers)
	mux.HandleFunc("GET /locations/{location_id}/cases", apiCfg.handlerCasesGetByLocation)
	mux.HandleFunc("GET /cases/{case_id}", apiCfg.handlerCaseGetByID)

	mux.HandleFunc("POST /locations/{location_id}/members", apiCfg.handlerAddLocationMember)
	mux.HandleFunc("DELETE /locations/{location_id}/members", apiCfg.handlerRemoveLocationMember)

	mux.HandleFunc("GET /search/users", apiCfg.handlerUsersGetByEmail)
	mux.HandleFunc("GET /search/locations/", apiCfg.handlerLocationsGetByOwner)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
