package handler

import (
	"encoding/json"
	"fmt"
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
		fmt.Println("Received POST request /api/shorturl")

		var req struct {
			URL string `json:"url"`
		}

		err := json.NewDecoder(r.Body).Decode(&req)

		fmt.Println("Request:", req)

		if err != nil {
			fmt.Println("Error decoding request:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		url, err := h.service.CreateUrl(userId, req.URL)
		if err != nil {
			fmt.Println("Error creating URL:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("Created URL:", url)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(url)
		return
	}

	if method == http.MethodGet {
		fmt.Println("Received GET request /api/shorturl")

		urls := h.service.ListUrlsByUserID(userId)
		fmt.Println("URLs:", urls)
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

	fmt.Println("Received GET request /api/s/:code")

	code := strings.TrimPrefix(r.URL.Path, "/api/s/")

	url, err := h.service.Resolve(r.Context(), code)
	if err != nil {
		fmt.Println("Error resolving URL:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if url == "" {
		fmt.Println("URL not found")
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
