package handler

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"shortener/internal/model"
	"shortener/internal/repository"
	"strings"

	"github.com/go-chi/chi/v5"
)

const shortCodeLength = 6

type URLHandler struct {
	repo    *repository.URLRepository
	baseURL string
}

func NewURLHandler(repo *repository.URLRepository, baseURL string) *URLHandler {
	return &URLHandler{
		repo:    repo,
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req model.UrlRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		RespondJson(w, r, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	raw := strings.TrimSpace(req.Url)
	if raw == "" {
		RespondJson(w, r, http.StatusBadRequest, map[string]string{"error": "url is required"})
		return
	}

	if _, err := validateHTTPURL(raw); err != nil {
		RespondJson(w, r, http.StatusBadRequest, map[string]string{"errror": "invalid URL; only http/https allowed"})
		return
	}

	for i := 0; i < 6; i++ {
		code, err := generateShortCode(shortCodeLength)
		if err != nil {
			RespondJson(w, r, http.StatusInternalServerError, map[string]string{"error": "failed to generate short code"})
			return
		}

		item, err := h.repo.Create(r.Context(), raw, code)
		if err != nil {
			if isUniqueShortCodeErr(err) {
				continue
			}
			RespondJson(w, r, http.StatusInternalServerError, map[string]string{"error": "failed to create short URL"})
			return
		}

		RespondJson(w, r, http.StatusCreated, model.ShortenResponse{
			Id:        item.Id,
			Url:       item.Url,
			ShortCode: item.ShortCode,
			ShortURL:  fmt.Sprintf("%s/%s", h.baseURL, item.ShortCode),
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
		return
	}

	RespondJson(w, r, http.StatusInternalServerError, map[string]string{"error": "failed to allocate unique short code"})
}

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := strings.TrimSpace(chi.URLParam(r, "shortCode"))
	if shortCode == "" {
		RespondJson(w, r, http.StatusBadRequest, map[string]string{"error": "short code is required"})
		return
	}

	item, err := h.repo.GetByShortCode(r.Context(), shortCode)
	if err != nil {
		if errors.Is(err, repository.ErrURLNotFound) {
			RespondJson(w, r, http.StatusNotFound, map[string]string{"error": "short URL not found"})
			return
		}
		RespondJson(w, r, http.StatusInternalServerError, map[string]string{"error": "failed to resolve short URL"})
		return
	}
	http.Redirect(w, r, item.Url, http.StatusFound)
}

func validateHTTPURL(raw string) (*url.URL, error) {
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("unsupported scheme")
	}
	if u.Host == "" {
		return nil, errors.New("missing host")
	}
	return u, nil
}

func generateShortCode(length int) (string, error) {
	const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	rb := make([]byte, length)

	if _, err := rand.Read(rb); err != nil {
		return "", err
	}

	for i := range b {
		b[i] = alphabet[int(rb[i])%len(alphabet)]
	}
	return string(b), nil
}

func isUniqueShortCodeErr(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed: urls.short_code")
}
