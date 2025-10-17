package main

import (
	"encoding/json"
	"net/http"

	"github.com/thachnguyensg/chirpy/internal/auth"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	responseWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Email:     user.Email,
	})
}
