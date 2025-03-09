package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rodabaugh/digitalshelf/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerBooksGet(w http.ResponseWriter, r *http.Request) {
	dbBooks, err := cfg.db.GetBooks(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get books from database", err)
		return
	}

	books := []Book{}

	for _, dbBook := range dbBooks {
		books = append(books, Book{
			ID:              dbBook.ID,
			Title:           dbBook.Title,
			Author:          dbBook.Author,
			Genre:           dbBook.Genre,
			Barcode:         dbBook.Barcode,
			ShelfID:         dbBook.ShelfID,
			PublicationDate: dbBook.PublicationDate,
			CreatedAt:       dbBook.CreatedAt,
			UpdatedAt:       dbBook.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, books)
}

func (cfg *apiConfig) handlerBooksGetByShelf(w http.ResponseWriter, r *http.Request) {
	shelfIDString := r.PathValue("shelf_id")
	if shelfIDString == "" {
		respondWithError(w, http.StatusBadRequest, "No shelf id was provided", fmt.Errorf("no shelf_id was provided"))
		return
	}

	shelfID, err := uuid.Parse(shelfIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid shelf ID", err)
		return
	}

	// Validate user is authorized to get books at the location of shelf.
	shelfLocation, err := cfg.db.GetShelfLocation(r.Context(), shelfID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get shelf location", err)
		return
	}

	err = cfg.authorizeMember(shelfLocation.ID, *r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to get books at the location of that shelf", err)
		return
	}

	dbBooks, err := cfg.db.GetBooksByShelf(r.Context(), shelfID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No books found for that shelf", err)
		return
	}

	books := []Book{}

	for _, dbBook := range dbBooks {
		books = append(books, Book{
			ID:              dbBook.ID,
			Title:           dbBook.Title,
			Author:          dbBook.Author,
			Genre:           dbBook.Genre,
			Barcode:         dbBook.Barcode,
			ShelfID:         dbBook.ShelfID,
			PublicationDate: dbBook.PublicationDate,
			CreatedAt:       dbBook.CreatedAt,
			UpdatedAt:       dbBook.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, books)
}

func (cfg *apiConfig) handlerBookGetByID(w http.ResponseWriter, r *http.Request) {
	bookIDString := r.PathValue("book_id")
	if bookIDString == "" {
		respondWithError(w, http.StatusBadRequest, "No book id was provided", fmt.Errorf("no book id was provided"))
		return
	}

	bookID, err := uuid.Parse(bookIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid book ID", err)
		return
	}

	// Validate user is authorized to get books at the location of requested book.
	bookLocation, err := cfg.db.GetBookLocation(r.Context(), bookID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get book location", err)
		return
	}

	err = cfg.authorizeMember(bookLocation.ID, *r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to get books at this location", err)
		return
	}

	dbBook, err := cfg.db.GetBookByID(r.Context(), bookID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Book not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Book{
		ID:              dbBook.ID,
		Title:           dbBook.Title,
		Author:          dbBook.Author,
		Genre:           dbBook.Genre,
		Barcode:         dbBook.Barcode,
		ShelfID:         dbBook.ShelfID,
		PublicationDate: dbBook.PublicationDate,
		CreatedAt:       dbBook.CreatedAt,
		UpdatedAt:       dbBook.UpdatedAt,
	})
}

func (cfg *apiConfig) handlerGetBookByBarcode(w http.ResponseWriter, r *http.Request) {
	barcode := r.PathValue("barcode")
	fmt.Println(barcode)
	if barcode == "" {
		respondWithError(w, http.StatusBadRequest, "No barcode was provided", fmt.Errorf("no barcode was provided"))
		return
	}

	dbBook, err := cfg.db.GetBookByBarcode(r.Context(), barcode)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Book not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Book{
		ID:              dbBook.ID,
		Title:           dbBook.Title,
		Author:          dbBook.Author,
		Genre:           dbBook.Genre,
		Barcode:         dbBook.Barcode,
		ShelfID:         dbBook.ShelfID,
		PublicationDate: dbBook.PublicationDate,
		CreatedAt:       dbBook.CreatedAt,
		UpdatedAt:       dbBook.UpdatedAt,
	})
}

func (cfg *apiConfig) handlerBooksGetByLocation(w http.ResponseWriter, r *http.Request) {
	locationIDString := r.PathValue("location_id")
	if locationIDString == "" {
		respondWithError(w, http.StatusBadRequest, "No location id was provided", fmt.Errorf("no location id was provided"))
		return
	}

	locationID, err := uuid.Parse(locationIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid location ID", err)
		return
	}

	// Validate user is permitted to get books for the location.
	err = cfg.authorizeMember(locationID, *r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to get books for this location", err)
		return
	}

	dbBooks, err := cfg.db.GetBooksByLocation(r.Context(), locationID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No books found for that location", err)
		return
	}

	books := []Book{}

	for _, dbBook := range dbBooks {
		books = append(books, Book{
			ID:              dbBook.ID,
			Title:           dbBook.Title,
			Author:          dbBook.Author,
			Genre:           dbBook.Genre,
			Barcode:         dbBook.Barcode,
			ShelfID:         dbBook.ShelfID,
			PublicationDate: dbBook.PublicationDate,
			CreatedAt:       dbBook.CreatedAt,
			UpdatedAt:       dbBook.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, books)
}

func (cfg *apiConfig) handlerSearchBooks(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		LocationID string `json:"location_id"`
		Query      string `json:"query"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	locationIDString := requestBody.LocationID
	if locationIDString == "" {
		respondWithError(w, http.StatusBadRequest, "No location id was provided", fmt.Errorf("no location id was provided"))
		return
	}

	locationID, err := uuid.Parse(locationIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid location ID", err)
		return
	}

	// Validate user is permitted to search books for the location.
	err = cfg.authorizeMember(locationID, *r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not authorized to search books for this location", err)
		return
	}

	query := requestBody.Query
	if query == "" {
		respondWithError(w, http.StatusBadRequest, "No search query was provided", fmt.Errorf("no search query was provided"))
		return
	}

	dbBooks, err := cfg.db.SearchBooks(r.Context(), database.SearchBooksParams{
		WebsearchToTsquery: query,
		ID:                 locationID,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No books found for that location", err)
		return
	}

	books := []Book{}

	for _, dbBook := range dbBooks {
		books = append(books, Book{
			ID:              dbBook.ID,
			Title:           dbBook.Title,
			Author:          dbBook.Author,
			Genre:           dbBook.Genre,
			Barcode:         dbBook.Barcode,
			ShelfID:         dbBook.ShelfID,
			PublicationDate: dbBook.PublicationDate,
			CreatedAt:       dbBook.CreatedAt,
			UpdatedAt:       dbBook.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, books)
}
