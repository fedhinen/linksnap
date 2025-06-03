package handler

import (
	"encoding/json"
	"linksnap/internal/auth"
	"linksnap/internal/service"
	"net/http"
	"strings"
)

type ShortUrlHandler struct {
	service *service.ShortUrlService
}

func NewURLHandler(service *service.ShortUrlService) *ShortUrlHandler {
	return &ShortUrlHandler{service: service}
}

func (h *ShortUrlHandler) ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	userId, err := auth.GetAuthenticatedUserID(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if userId == "" {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}

	if method == http.MethodPost {
		var req struct {
			URL string `json:"url"`
		}

		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		url, err := h.service.CreateUrl(userId, req.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(url)
		return
	}

	if method == http.MethodGet {
		urls := h.service.ListUrlsByUserID(userId)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(urls)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	return
}

func (h *ShortUrlHandler) GetShortURL(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	if method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/api/s/")

	url, err := h.service.Resolve(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if url == "" {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	type Response struct {
		Code string `json:"code"`
		URL  string `json:"url"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Code: code,
		URL:  url,
	})
}
