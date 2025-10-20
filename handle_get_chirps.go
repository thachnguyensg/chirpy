package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/thachnguyensg/chirpy/internal/database"
)

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	var returnChirps []database.Chirp
	var err error
	sortChirp := "asc"

	authorId := r.URL.Query().Get("author_id")
	sortChirp = r.URL.Query().Get("sort")

	if authorId != "" {
		userID, err := uuid.Parse(authorId)
		if err != nil {
			responseWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}

		returnChirps, err = cfg.db.GetChirpsByUserID(r.Context(), userID)
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, "Something went wrong", err)
			return
		}
	} else {
		returnChirps, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, "Something went wrong", err)
			return
		}
	}

	if sortChirp == "desc" {
		sort.Slice(returnChirps, func(i, j int) bool {
			return i > j
		})
	} else {
		sort.Slice(returnChirps, func(i, j int) bool {
			return i < j
		})
	}

	var response []Chirp
	for _, c := range returnChirps {
		response = append(response, Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}

	responseWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	cid := r.PathValue("chirp_id")
	chirpID, err := uuid.Parse(cid)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		responseWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	responseWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
