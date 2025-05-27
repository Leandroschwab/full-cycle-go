package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Leandroschwab/full-cycle-go/RateLimiter/internal/config"
	"github.com/Leandroschwab/full-cycle-go/RateLimiter/internal/limiter"
	"github.com/Leandroschwab/full-cycle-go/RateLimiter/internal/middleware"
	"github.com/Leandroschwab/full-cycle-go/RateLimiter/internal/storage"
)

func main() {
	// Carrega as configurações
	cfg := config.LoadConfig()

	// Inicializa o armazenamento Redis
	redisStorage := storage.NewRedisStorage(cfg.RedisURL)

	// Inicializa o serviço de rate limiting
	rateLimiter := limiter.NewService(redisStorage, cfg)
	defer rateLimiter.Close()

	// Inicializa o middleware
	rateLimiterMiddleware := middleware.NewRateLimiterMiddleware(rateLimiter)

	// Configura o router
	r := mux.NewRouter()

	// Adiciona o middleware ao router
	r.Use(rateLimiterMiddleware.Middleware)

	// Rota de exemplo
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Rate Limited World!"))
	})

	// Inicia o servidor
	log.Println("Server starting on :8080...")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
