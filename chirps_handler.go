package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thachnguyensg/chirpy/internal/database"
)

func (cfg *apiConfig) chirpsHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type returnedChirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	type response struct {
		Chirp returnedChirp `json:"chirp"`
	}

	var params parameters
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	validateChirp(w, params.Body)

	cleanedBody, err := cleanupInputV2(params.Body)
	if err != nil {
		fmt.Println("Error cleaning up input:", err)
		responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: params.UserID,
	})
	if err != nil {
		fmt.Println("Error creating chirp:", err)
		responseWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	responseWithJSON(w, http.StatusCreated, returnedChirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func validateChirp(w http.ResponseWriter, chirp string) {
	if len(chirp) > 140 {
		fmt.Printf("Chirp body exceeds 140 characters: %d\n", len(chirp))
		responseWithError(w, http.StatusBadRequest, "Chirp body exceeds 140 characters")
		return
	}
}

func cleanupInputV2(input string) (string, error) {
	fmt.Printf("input: %v\n", input)
	profanities := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	maskStr := "****"
	cleaned := strings.Builder{}
	wb := 0
	wl := 0
	for i := 0; i < len(input); i++ {
		if input[i] == ' ' {
			if wl > 0 {
				word := input[wb : wb+wl]
				if _, ok := profanities[strings.ToLower(word)]; ok {
					cleaned.WriteString(maskStr)
				} else {
					cleaned.WriteString(word)
				}
				fmt.Printf("%v ", word)

				wl = 0
			}
			cleaned.WriteString(string(input[i]))
		} else {
			if wl == 0 {
				wb = i
			}
			wl += 1
		}
	}
	if wl > 0 {
		fmt.Printf("%v \n", input[wb:wb+wl])
		word := input[wb : wb+wl]
		if _, ok := profanities[strings.ToLower(word)]; ok {
			cleaned.WriteString(maskStr)
		} else {
			cleaned.WriteString(word)
		}
	}
	fmt.Printf("cleaned input: %v\n", cleaned.String())

	return cleaned.String(), nil
}
