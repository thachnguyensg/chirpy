package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string
	}
	type valid struct {
		Valid       bool   `json:"valid,omitempty"`
		CleanedBody string `json:"cleaned_body,omitempty"`
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Printf("Error decoding parameters: %v\n", err)
		responseWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		fmt.Printf("Chirp body exceeds 140 characters: %d\n", len(params.Body))
		responseWithError(w, http.StatusBadRequest, "Chirp body exceeds 140 characters")
		return
	}

	cleanedBody, err := cleanupInputV2(params.Body)
	if err != nil {
		fmt.Printf("Error cleaning up input: %v\n", err)
		responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	responseWithJSON(w, http.StatusOK, valid{CleanedBody: cleanedBody})
}

func cleanupInput(input string) (string, error) {
	profanities := []string{"kerfuffle", "sharbert", "fornax"}
	maskStr := "****"
	re, err := regexp.Compile(`(?i)` + strings.Join(profanities, "|"))
	if err != nil {
		return "", fmt.Errorf("error compiling regex for profanity %v: %v", profanities, err)
	}
	input = re.ReplaceAllString(input, maskStr)
	return input, nil
}
