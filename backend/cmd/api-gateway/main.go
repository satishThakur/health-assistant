package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting API Gateway...")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "API Gateway is healthy")
	})

	port := ":8080"
	log.Printf("API Gateway listening on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
