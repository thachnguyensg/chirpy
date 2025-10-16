package main

import (
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) reset() {
	cfg.fileserverHits.Store(0)
}

func main() {
	const rootPath = "."
	const port = "8080"

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
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

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Server running at :%s\n", port)
	log.Fatal(server.ListenAndServe())
}
