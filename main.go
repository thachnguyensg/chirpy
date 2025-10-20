package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/thachnguyensg/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secretKey      string
}

func (cfg *apiConfig) reset() {
	cfg.fileserverHits.Store(0)
}

func main() {
	const rootPath = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secretKey := os.Getenv("SECRET_KEY")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:             database.New(db),
		platform:       platform,
		secretKey:      secretKey,
	}

	mux := http.NewServeMux()
	rp := http.Dir(rootPath)
	pwd, _ := os.Getwd()
	log.Printf("rootpath: %s\n", pwd)
	fs := http.FileServer(rp)
	mux.Handle("GET /app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", fs)))

	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)
	mux.HandleFunc("PUT /api/users", apiCfg.updateUserHandler)

	mux.HandleFunc("POST /api/login", apiCfg.loginHandler)
	mux.HandleFunc("POST /api/refresh", apiCfg.refreshTokenHandler)
	mux.HandleFunc("POST /api/revoke", apiCfg.revokeTokenHandler)

	mux.HandleFunc("GET /api/chirps", apiCfg.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirp_id}", apiCfg.getChirpHandler)
	mux.HandleFunc("DELETE /api/chirps/{chirp_id}", apiCfg.deleteChirpHandler)

	mux.Handle("POST /api/chirps", apiCfg.createChirpHandler())
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.userChirpyRedWebhookHandler)

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Server running at :%s\n", port)
	log.Fatal(server.ListenAndServe())
}
