package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type RequestBody struct {
	Url string `json:"url"`
}

type Url struct {
	Id        uuid.UUID `json:"id"`
	Url       string    `json:"url"`
	ShortCode string    `json:"shortCode"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

var urls map[uuid.UUID]Url = make(map[uuid.UUID]Url)
var urlsMu sync.RWMutex

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/shorten", CreateUrl)
	r.Get("/shorten/{shortCode}", RetrieveOriginalUrlFromShortUrl)
	r.Put("/shorten/{shortCode}", UpdateUrl)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func CreateUrl(w http.ResponseWriter, r *http.Request) {
	var body RequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return 
	}

	if body.Url == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return 
	}

	id := uuid.New()
	now := time.Now().UTC()
	created := Url{ 
		Id:        id,
		Url:       body.Url,
		ShortCode: id.String()[:6],
		CreatedAt: now,
		UpdatedAt: now,
	}

	urlsMu.Lock()
	urls[id] = created
	urlsMu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(created); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func RetrieveOriginalUrlFromShortUrl(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode")
	var found Url
	var ok bool

	urlsMu.RLock()
	for _, u := range urls {
		if u.ShortCode == shortCode {
			found = u
			ok = true
			break
		}
	}
	urlsMu.RUnlock()

	if !ok {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(found); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func UpdateUrl(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode")	
	var body RequestBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return 
	}

	if body.Url == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return 
	}

	var (
		targetID uuid.UUID
		found    Url
		ok       bool
	)

	urlsMu.Lock()
	for id, u := range urls {
		if u.ShortCode == shortCode {
			targetID = id
			found = u
			ok = true
			break
		}
	}

	if !ok {
		urlsMu.Unlock()
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	found.Url = body.Url
	found.UpdatedAt = time.Now().UTC()
	urls[targetID] = found
	urlsMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(found); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
