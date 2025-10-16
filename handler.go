package main

import (
	"fmt"
	"net/http"
)

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	payload := fmt.Sprintf(`
			<html>
			  <body>
			    <h1>Welcome, Chirpy Admin</h1>
			    <p>Chirpy has been visited %d times!</p>
			  </body>
			</html>
			`, cfg.fileserverHits.Load())

	w.Write([]byte(payload))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		fmt.Printf("Reset attempted on non-dev platform: %s\n", cfg.platform)
		responseWithError(w, http.StatusForbidden, "Forbidden")
		return
	}

	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		fmt.Printf("Error deleting all users: %v\n", err)
		responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	responseWithJSON(w, http.StatusOK, map[string]string{"message": "All users deleted"})

	// cfg.reset()
	// hits := cfg.fileserverHits.Load()
	// w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte("Hits: " + fmt.Sprintf("%d", hits)))
}
