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
	platform  string
	db        *database.Queries
	jwtSecret string
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
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
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

	mux.HandleFunc("GET /admin/healthz", readinessEndpoint)

	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUsersUpdate)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	mux.HandleFunc("POST /api/revoke-all", apiCfg.handlerRevokeSessions)

	mux.HandleFunc("POST /api/locations", apiCfg.handlerLocationsCreate)
	mux.HandleFunc("POST /api/cases", apiCfg.handlerCasesCreate)
	mux.HandleFunc("GET /api/cases", apiCfg.handlerCaseGet)
	mux.HandleFunc("POST /api/shelves", apiCfg.handlerShelfCreate)
	mux.HandleFunc("GET /api/shelves", apiCfg.handlerShelvesGet)
	mux.HandleFunc("POST /api/movies", apiCfg.handlerMovieCreate)
	mux.HandleFunc("GET /api/movies", apiCfg.handlerMoviesGet)
	mux.HandleFunc("POST /api/shows", apiCfg.handlerShowCreate)
	mux.HandleFunc("GET /api/shows", apiCfg.handlerShowsGet)

	mux.HandleFunc("GET /api/users", apiCfg.handlerUsersGet)
	mux.HandleFunc("GET /api/users/{user_id}", apiCfg.handlerUserGetByID)
	mux.HandleFunc("GET /api/users/{user_id}/locations", apiCfg.handlerGetUserLocations)
	mux.HandleFunc("GET /api/users/{user_id}/invites", apiCfg.handlerGetUserInvites)
	mux.HandleFunc("GET /api/locations", apiCfg.handlerLocationsGet)
	mux.HandleFunc("GET /api/locations/{location_id}", apiCfg.handlerLocationsGetByID)
	mux.HandleFunc("GET /api/locations/{location_id}/members", apiCfg.handlerGetLocationMembers)
	mux.HandleFunc("GET /api/locations/{location_id}/invites", apiCfg.handlerGetLocationInvites)
	mux.HandleFunc("GET /api/locations/{location_id}/cases", apiCfg.handlerCasesGetByLocation)
	mux.HandleFunc("GET /api/locations/{location_id}/movies", apiCfg.handlerMoviesGetByLocation)
	mux.HandleFunc("GET /api/locations/{location_id}/shows", apiCfg.handlerShowsGetByLocation)
	mux.HandleFunc("GET /api/cases/{case_id}", apiCfg.handlerCaseGetByID)
	mux.HandleFunc("GET /api/cases/{case_id}/shelves", apiCfg.handlerShelvesGetByCase)
	mux.HandleFunc("GET /api/shelves/{shelf_id}", apiCfg.handlerShelfGetByID)
	mux.HandleFunc("GET /api/shelves/{shelf_id}/movies", apiCfg.handlerMoviesGetByShelf)
	mux.HandleFunc("GET /api/movies/{movie_id}", apiCfg.handlerMovieGetByID)
	mux.HandleFunc("GET /api/shelves/{shelf_id}/shows", apiCfg.handlerShowsGetByShelf)
	mux.HandleFunc("GET /api/shows/{show_id}", apiCfg.handlerShowGetByID)

	mux.HandleFunc("DELETE /api/locations/{location_id}/members/{user_id}", apiCfg.handlerRemoveLocationMember)
	mux.HandleFunc("POST /api/locations/{location_id}/members", apiCfg.handlerAddLocationMember)
	mux.HandleFunc("DELETE /api/locations/{location_id}/invites/{user_id}", apiCfg.handlerRemoveLocationInvite)
	mux.HandleFunc("POST /api/locations/{location_id}/invites", apiCfg.handlerAddLocationInvite)

	mux.HandleFunc("GET /api/search/users", apiCfg.handlerUsersGetByEmail)
	mux.HandleFunc("GET /api/search/locations/", apiCfg.handlerLocationsGetByOwner)
	mux.HandleFunc("GET /api/search/movie_barcodes/{barcode}", apiCfg.handlerGetMovieByBarcode)
	mux.HandleFunc("GET /api/search/movies", apiCfg.handlerSearchMovies)
	mux.HandleFunc("GET /api/search/show_barcodes/{barcode}", apiCfg.handlerGetShowByBarcode)
	mux.HandleFunc("GET /api/search/shows", apiCfg.handlerSearchShows)

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
