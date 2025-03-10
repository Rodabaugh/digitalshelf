package main

import (
	"database/sql"
	"fmt"
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
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Using environment variables.")
	} else {
		fmt.Println("Loaded .env file.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

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
	mux.HandleFunc("PUT /api/movies/{movie_id}", apiCfg.handlerMoviesUpdate)
	mux.HandleFunc("POST /api/shows", apiCfg.handlerShowCreate)
	mux.HandleFunc("GET /api/shows", apiCfg.handlerShowsGet)
	mux.HandleFunc("POST /api/books", apiCfg.handlerBookCreate)
	mux.HandleFunc("GET /api/books", apiCfg.handlerBooksGet)
	mux.HandleFunc("POST /api/music", apiCfg.handlerMusicCreate)
	mux.HandleFunc("GET /api/music", apiCfg.handlerMusicGet)

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
	mux.HandleFunc("GET /api/locations/{location_id}/books", apiCfg.handlerBooksGetByLocation)
	mux.HandleFunc("GET /api/locations/{location_id}/music", apiCfg.handlerMusicGetByLocation)
	mux.HandleFunc("GET /api/cases/{case_id}", apiCfg.handlerCaseGetByID)
	mux.HandleFunc("GET /api/cases/{case_id}/shelves", apiCfg.handlerShelvesGetByCase)
	mux.HandleFunc("GET /api/shelves/{shelf_id}", apiCfg.handlerShelfGetByID)
	mux.HandleFunc("GET /api/shelves/{shelf_id}/movies", apiCfg.handlerMoviesGetByShelf)
	mux.HandleFunc("GET /api/movies/{movie_id}", apiCfg.handlerMovieGetByID)
	mux.HandleFunc("GET /api/shelves/{shelf_id}/shows", apiCfg.handlerShowsGetByShelf)
	mux.HandleFunc("GET /api/shows/{show_id}", apiCfg.handlerShowGetByID)
	mux.HandleFunc("GET /api/shelves/{shelf_id}/books", apiCfg.handlerBooksGetByShelf)
	mux.HandleFunc("GET /api/books/{book_id}", apiCfg.handlerBookGetByID)
	mux.HandleFunc("GET /api/shelves/{shelf_id}/music", apiCfg.handlerMusicGetByShelf)
	mux.HandleFunc("GET /api/music/{music_id}", apiCfg.handlerMusicGetByID)

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
	mux.HandleFunc("GET /api/search/book_barcodes/{barcode}", apiCfg.handlerGetBookByBarcode)
	mux.HandleFunc("GET /api/search/books", apiCfg.handlerSearchBooks)
	mux.HandleFunc("GET /api/search/music_barcodes/{barcode}", apiCfg.handlerGetMusicByBarcode)
	mux.HandleFunc("GET /api/search/music", apiCfg.handlerSearchMusic)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Serving DigitalShelf backend on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
