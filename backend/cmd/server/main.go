package main

import (
	"log"
	"net/http"
	"os"

	"calculator-backend/internal/httpapi"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("calculator backend listening on :%s", port)
	if err := http.ListenAndServe(":"+port, httpapi.NewRouter()); err != nil {
		log.Fatal(err)
	}
}
