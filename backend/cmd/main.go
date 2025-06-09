package main

import (
	"linksnap/internal/config"
	"linksnap/internal/handler"
	"linksnap/internal/middleware"
	"linksnap/internal/service"
	"linksnap/internal/storage"
	"log"
	"net/http"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkHttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, continuing without it")
	}

	envs := config.LoadEnv()

	clerk.SetKey(envs.ClerkSecretKey)

	db, err := storage.InitializeDatabase(envs)
	if err != nil {
		log.Fatal("DB init error:", err)
	}
	defer db.Close()

	// Inicializar Redis (Valkey)
	redis, err := storage.InitializeRedis(envs)
	if err != nil {
		log.Fatal("Redis init error:", err)
	}

	var store storage.UrlStore

	if envs.DatabaseDriver == "sqlite3" {
		store = storage.NewSqliteStore(db)
	} else if envs.DatabaseDriver == "postgres" {
		store = storage.NewPostgresStore(db)
	}

	urlService := service.NewShortUrlService(store, *redis)
	urlHandler := handler.NewURLHandler(urlService)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/s/", urlHandler.GetShortURLHandler)
	mux.HandleFunc("/api/health/", handler.HealthHandler)

	protectedRoute := http.HandlerFunc(urlHandler.ShortURLHandler)
	mux.Handle("/api/shorturl/", clerkHttp.WithHeaderAuthorization()(protectedRoute))

	handlerWithCORS := middleware.CORS(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5271"
	}

	log.Printf("Listening on port %s...", port)
	err = http.ListenAndServe(":"+port, handlerWithCORS)
	if err != nil {
		log.Fatal(err)
	}
}
