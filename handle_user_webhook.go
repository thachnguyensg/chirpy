package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/thachnguyensg/chirpy/internal/auth"
)

func (cfg *apiConfig) userChirpyRedWebhookHandler(w http.ResponseWriter, r *http.Request) {
	_, err := auth.GetApiKey(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	var parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&parameters)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if parameters.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.GetUserByID(r.Context(), parameters.Data.UserID)
	if err != nil {
		responseWithError(w, http.StatusNotFound, "Failed to fetch user", err)
		return
	}

	_, err = cfg.db.SetUserChirpyRed(r.Context(), parameters.Data.UserID)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
