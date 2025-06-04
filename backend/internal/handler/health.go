package handler

import (
	"encoding/json"
	"fmt"
	"linksnap/internal/service"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	if method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := service.GetHealth()

	fmt.Println("Health response:", response)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
