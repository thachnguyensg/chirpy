package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/thachnguyensg/chirpy/internal/auth"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int64  `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Something went wrong", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	errorMsg := "invalid email or password"
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, errorMsg, err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		responseWithError(w, http.StatusUnauthorized, errorMsg, err)
		return
	}

	tokenExpiry := time.Hour * 1
	if params.ExpiresInSeconds > 0 {
		tokenExpiry = time.Second * time.Duration(params.ExpiresInSeconds)
	}
	token, err := auth.MakeJWT(user.ID, cfg.secretKey, tokenExpiry)

	responseWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Email:     user.Email,
		Token:     token,
	})
}

func getAuthenticatedUserID(header http.Header, secret string) (uuid.UUID, error) {
	token, err := auth.GetBearerToken(header)
	if err != nil {
		return uuid.Nil, err
	}
	userID, err := auth.ValidateJWT(token, secret)
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}
