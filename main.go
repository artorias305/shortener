package main

import (
	"log"
	"net/http"
	"os"

	"shortener/internal/api/handler"
	"shortener/internal/api/router"
	"shortener/internal/db"
	"shortener/internal/repository"
)

func main() {
	database, err := db.InitSQLite("data/shortener.db")
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer database.Close()

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	urlRepo := repository.NewURLRepository(database)
	urlHandler := handler.NewURLHandler(urlRepo, baseURL)

	r := router.New(urlHandler)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
