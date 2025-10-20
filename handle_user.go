package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/thachnguyensg/chirpy/internal/auth"
	"github.com/thachnguyensg/chirpy/internal/database"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Printf("Error decoding parameters: %v\n", err)
		responseWithError(w, http.StatusBadRequest, "Something went wrong", err)
		return
	}

	// err = validateUserPassword(params.Password)
	// if err != nil {
	// 	responseWithError(w, http.StatusBadRequest, err.Error(), err)
	// 	return
	// }

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	responseWithJSON(w, http.StatusCreated, User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt.Time,
		UpdatedAt:   user.UpdatedAt.Time,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func validateUserPassword(pw string) error {
	if len(pw) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	return nil
}

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	var parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	uid, err := getAuthenticatedUserID(r.Header, cfg.secretKey)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&parameters)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Cannot read email, password", err)
		return
	}

	hashedPassword, err := auth.HashPassword(parameters.Password)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Cannot hash password", err)
		return
	}

	user, err := cfg.db.UpdateUserAuth(r.Context(), database.UpdateUserAuthParams{
		ID:             uid,
		Email:          parameters.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Cannot update user", err)
		return
	}

	responseWithJSON(w, http.StatusOK, User{
		ID:          user.ID,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt.Time,
		UpdatedAt:   user.UpdatedAt.Time,
		IsChirpyRed: user.IsChirpyRed,
	})
}
