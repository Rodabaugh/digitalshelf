package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte("Reset is only allowed in dev environment."))
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to reset database", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Database reset to initial state."))
	if err != nil {
		fmt.Println(err)
	}
}
