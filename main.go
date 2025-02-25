package main

import (
	"log"
	"net/http"
)

func main() {
	const rootPath = "."
	const port = "8080"

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(rootPath))
	mux.Handle("/app/", http.StripPrefix("/app/", fs))

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Server running at :%s\n", port)
	log.Fatal(server.ListenAndServe())
}
