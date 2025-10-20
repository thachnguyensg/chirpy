package main

import (
	"net/http"
	"time"

	"github.com/thachnguyensg/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	newToken, err := auth.MakeJWT(user.ID, cfg.secretKey, 1*time.Hour)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Can not create new token", err)
		return
	}

	responseWithJSON(w, http.StatusOK, map[string]string{
		"token": newToken,
	})
}

func (cfg *apiConfig) revokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	_, err = cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Can not revoke token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
