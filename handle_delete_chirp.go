package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirp_id")
	chirpUID, err := uuid.Parse(chirpID)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	uid, err := getAuthenticatedUserID(r.Header, cfg.secretKey)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Cannot authenticate user", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpUID)
	if err != nil {
		responseWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}
	if chirp.UserID != uid {
		responseWithError(w, http.StatusForbidden, "Cannot delete chirp of another user", nil)
		return
	}

	_, err = cfg.db.DeleteChirp(r.Context(), chirpUID)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Cannot delete chirp", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
