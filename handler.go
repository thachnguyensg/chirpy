package main

import (
	"encoding/json"
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

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, _ *http.Request) {
	cfg.reset()
	hits := cfg.fileserverHits.Load()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits: " + fmt.Sprintf("%d", hits)))
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string
	}
	type response struct {
		Valid bool   `json:"valid,omitempty"`
		Error string `json:"error,omitempty"`
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Printf("Error decoding parameters: %v\n", err)
		resp := response{Error: "Something went wrong"}
		payload, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(payload)
		return
	}

	if len(params.Body) > 140 {
		fmt.Printf("Chirp body exceeds 140 characters: %d\n", len(params.Body))
		payload, _ := json.Marshal(response{Valid: false, Error: "Chirp is too long"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(payload)
		return
	}

	payload, _ := json.Marshal(response{Valid: true})
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}
