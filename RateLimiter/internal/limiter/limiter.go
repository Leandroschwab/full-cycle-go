package limiter

import (
	"context"
	"time"

	"github.com/Leandroschwab/full-cycle-go/RateLimiter/internal/config"
	"github.com/Leandroschwab/full-cycle-go/RateLimiter/internal/storage"
)

// RateLimitInfo contains information about the rate limit status
type RateLimitInfo struct {
	Allowed      bool
	CurrentCount int
	Limit        int
	Key          string
}

// RateLimiter define a interface para o serviço de rate limiting
type RateLimiter interface {
	// Verifica se uma requisição pode ser processada
	Allow(ctx context.Context, ip, token string) (RateLimitInfo, error)
}

// Service implementa a interface RateLimiter
type Service struct {
	storage       storage.Storage
	ipLimit       int
	tokenLimit    int
	blockDuration time.Duration
}

// NewService cria uma nova instância do serviço de rate limiting
func NewService(store storage.Storage, cfg *config.Config) *Service {
	return &Service{
		storage:       store,
		ipLimit:       cfg.IPRateLimit,
		tokenLimit:    cfg.TokenRateLimit,
		blockDuration: time.Duration(cfg.BlockDuration) * time.Second,
	}
}

// Allow verifica se uma requisição pode ser processada
func (s *Service) Allow(ctx context.Context, ip, token string) (RateLimitInfo, error) {
	// Se o token estiver presente, verifica primeiro o token
	if token != "" {
		return s.checkLimit(ctx, "token:"+token, s.tokenLimit)
	}

	// Se não houver token, verifica por IP
	return s.checkLimit(ctx, "ip:"+ip, s.ipLimit)
}

// checkLimit verifica se uma chave atingiu seu limite de requisições
func (s *Service) checkLimit(ctx context.Context, key string, limit int) (RateLimitInfo, error) {
	info := RateLimitInfo{
		Key:   key,
		Limit: limit,
	}

	// Verifica se a chave está bloqueada
	blocked, err := s.storage.IsBlocked(ctx, key)
	if err != nil {
		return info, err
	}

	if blocked {
		info.Allowed = false
		info.CurrentCount = limit + 1
		return info, nil
	}

	// Incrementa o contador para a chave
	count, err := s.storage.Increment(ctx, key)
	if err != nil {
		return info, err
	}

	info.CurrentCount = count
	info.Allowed = count <= limit

	// Se excedeu o limite, bloqueia a chave
	if !info.Allowed {
		err = s.storage.Block(ctx, key, s.blockDuration)
		if err != nil {
			return info, err
		}
	}

	return info, nil
}

func (s *Service) Close() error {
	return s.storage.Close()
}
