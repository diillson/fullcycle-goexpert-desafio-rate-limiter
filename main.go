package main

import (
	"fullcycle-goexpert-desafio-rate-limiter/limiter"
	"fullcycle-goexpert-desafio-rate-limiter/middleware"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379/0" // default value
	}

	// Inicialize o Redis storage
	storage := limiter.NewRedisStorage(redisURL)
	defer storage.Close()

	// Configure o rate limiter
	rateLimiter := limiter.NewRateLimiter(storage)

	// Crie o middleware
	limitMiddleware := middleware.NewRateLimitMiddleware(rateLimiter)

	// Rota de exemplo
	http.HandleFunc("/api", limitMiddleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}))

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
