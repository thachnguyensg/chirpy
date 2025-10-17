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
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token,omitempty"`
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
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Email:     user.Email,
	})
}

func validateUserPassword(pw string) error {
	if len(pw) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	return nil
}
