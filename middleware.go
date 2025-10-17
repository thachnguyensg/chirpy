package main

import (
	"net/http"

	"github.com/thachnguyensg/chirpy/internal/auth"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msg := "Unauthorized"
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			responseWithError(w, http.StatusUnauthorized, msg, err)
			return
		}
		_, err = auth.ValidateJWT(token, cfg.secretKey)
		if err != nil {
			responseWithError(w, http.StatusUnauthorized, msg, err)
			return
		}
		next.ServeHTTP(w, r)
	})
}
