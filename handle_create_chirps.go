package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thachnguyensg/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirpHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Body string `json:"body"`
		}

		userID, err := getAuthenticatedUserID(r.Header, cfg.secretKey)
		if err != nil {
			responseWithError(w, http.StatusUnauthorized, "Unauthorized", err)
			return
		}

		var params parameters
		err = json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			responseWithError(w, http.StatusBadRequest, err.Error(), err)
			return
		}

		cleanedChirp, err := validateChirp(params.Body)
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, err.Error(), err)
			return
		}

		chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
			Body:   cleanedChirp,
			UserID: userID,
		})
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, "Something went wrong", err)
			return
		}

		responseWithJSON(w, http.StatusCreated, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	})
}

func validateChirp(chirp string) (string, error) {
	if len(chirp) > 140 {
		return "", errors.New("Chirp body exceeds 140 characters")
	}
	cleanedChirp, err := cleanupInputV2(chirp)
	if err != nil {
		return "", err
	}
	return cleanedChirp, nil
}

func cleanupInputV2(input string) (string, error) {
	// fmt.Printf("input: %v\n", input)
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
				// fmt.Printf("%v ", word)

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
		// fmt.Printf("%v \n", input[wb:wb+wl])
		word := input[wb : wb+wl]
		if _, ok := profanities[strings.ToLower(word)]; ok {
			cleaned.WriteString(maskStr)
		} else {
			cleaned.WriteString(word)
		}
	}
	// fmt.Printf("cleaned input: %v\n", cleaned.String())

	return cleaned.String(), nil
}
