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
	fmt.Println("Received request /api/shorturl")

	method := r.Method

	fmt.Println("Method:", method)

	userId, err := auth.GetAuthenticatedUserID(r.Context())
	if err != nil {
		fmt.Println("Error getting authenticated user ID:", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if userId == "" {
		fmt.Println("User ID not found")
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/shorturl/")

	if id != "/api/shorturl/" && method == http.MethodDelete {
		fmt.Println("Received DELETE request /api/shorturl")

		deletedUrl, err := h.service.DeleteUrl(userId, id)
		if err != nil {
			fmt.Println("Error deleting URL:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("Deleted URL:", deletedUrl)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(deletedUrl)
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

func (h *ShortUrlHandler) GetShortURLHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received GET request /api/s/:code")

	method := r.Method

	fmt.Println("Method:", method)

	if method != http.MethodGet {
		fmt.Println("Method not allowed")
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
