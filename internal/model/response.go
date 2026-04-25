package model

import (
	"time"

	"github.com/google/uuid"
)

type UrlResponse struct {
	Id        uuid.UUID `json:"id"`
	Url       string    `json:"url"`
	ShortCode string    `json:"shortCode"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ShortenResponse struct {
	Id        uuid.UUID `json:"id"`
	Url       string    `json:"url"`
	ShortCode string    `json:"shortCode"`
	ShortURL  string    `json:"shortUrl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
