package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting Data Service...")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Data Service is healthy")
	})

	port := ":8081"
	log.Printf("Data Service listening on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
