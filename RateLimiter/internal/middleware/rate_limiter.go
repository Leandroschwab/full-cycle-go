package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/Leandroschwab/full-cycle-go/RateLimiter/internal/limiter"
)

const (
	// Cabeçalho para o token de API
	ApiKeyHeader = "API_KEY"
)

// RateLimiterMiddleware é um middleware para controlar o rate limiting
type RateLimiterMiddleware struct {
	limiter limiter.RateLimiter
}

// NewRateLimiterMiddleware cria uma nova instância do middleware
func NewRateLimiterMiddleware(limiter limiter.RateLimiter) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		limiter: limiter,
	}
}

// Middleware retorna o handler HTTP para integração com o servidor web
func (m *RateLimiterMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtém o IP do cliente
		ip := getClientIP(r)

		// Obtém o token de API do cabeçalho, se presente
		token := r.Header.Get(ApiKeyHeader)

		// Verifica se a requisição pode ser processada
		info, err := m.limiter.Allow(r.Context(), ip, token)
		if err != nil {
			log.Printf("Rate Limit Error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Identifica o tipo de limite (IP ou Token)
		limitType := "IP"
		identifier := ip
		if token != "" {
			limitType = "Token"
			identifier = token
		}

		// Log rate limit information to server console
		log.Printf("Rate Limit - Type: %s, Identifier: %s, Count: %d/%d, Remaining: %d, Allowed: %t",
			limitType, identifier, info.CurrentCount, info.Limit, info.Limit-info.CurrentCount, info.Allowed)

		if !info.Allowed {
			log.Printf("Rate Limit Exceeded - Type: %s, Identifier: %s", limitType, identifier)
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
			return
		}

		// Processa a requisição normalmente
		next.ServeHTTP(w, r)
	})
}

// getClientIP extrai o endereço IP do cliente da requisição
func getClientIP(r *http.Request) string {
	// Verifica se há um IP encaminhado por um proxy
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Se tiver múltiplos IPs, pega apenas o primeiro
		if commaIdx := strings.Index(forwarded, ","); commaIdx != -1 {
			forwarded = forwarded[:commaIdx]
		}
		return strings.TrimSpace(forwarded)
	}

	// Usa o IP de origem da requisição, removendo a porta
	remoteAddr := r.RemoteAddr
	if colonIdx := strings.LastIndex(remoteAddr, ":"); colonIdx != -1 {
		remoteAddr = remoteAddr[:colonIdx]
	}

	// Remove colchetes IPv6 se presentes
	remoteAddr = strings.TrimPrefix(remoteAddr, "[")
	remoteAddr = strings.TrimSuffix(remoteAddr, "]")

	return remoteAddr
}
