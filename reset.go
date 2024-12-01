package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}

	cfg.db.ResetUsers(r.Context())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Database reset to initial state."))
}
