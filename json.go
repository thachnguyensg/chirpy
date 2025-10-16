package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func responseWithError(w http.ResponseWriter, status int, message string) {
	type response struct {
		Error string `json:"error"`
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	payload, err := json.Marshal(response{Error: message})
	if err != nil {
		fmt.Printf("Error marshalling error response: %v\n", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	w.Write(payload)
}

func responseWithJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	p, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshalling response: %v\n", err)
		responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	w.Write(p)
}
