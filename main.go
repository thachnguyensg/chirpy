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
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", fs)))
	mux.Handle("/healthz", apiCfg.middlewareMetricsInc(http.HandlerFunc(healthzHandler)))
	mux.Handle("/metrics", metricsHandler(apiCfg))
	mux.Handle("/reset", resetHandler(apiCfg))

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Server running at :%s\n", port)
	log.Fatal(server.ListenAndServe())
}
