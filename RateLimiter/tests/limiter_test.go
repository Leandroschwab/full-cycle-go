package tests

import (
	"context"
	"testing"
	"time"

	"github.com/Leandroschwab/full-cycle-go/RateLimiter/internal/config"
	"github.com/Leandroschwab/full-cycle-go/RateLimiter/internal/limiter"
)

// MockStorage implementa a interface storage.Storage para testes
type MockStorage struct {
	counters    map[string]int
	blockedKeys map[string]bool
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		counters:    make(map[string]int),
		blockedKeys: make(map[string]bool),
	}
}

func (s *MockStorage) Increment(ctx context.Context, key string) (int, error) {
	s.counters[key]++
	return s.counters[key], nil
}

func (s *MockStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	return s.blockedKeys[key], nil
}

func (s *MockStorage) Block(ctx context.Context, key string, duration time.Duration) error {
	s.blockedKeys[key] = true
	// When a key is blocked, its counter should be at or above the limit
	// This ensures proper testing of the blocking mechanism
	if s.counters[key] < 1 {
		s.counters[key] = 1 // Ensure there's at least one count
	}
	return nil
}

func (s *MockStorage) Reset(ctx context.Context, key string) error {
	s.counters[key] = 0
	return nil
}

func (s *MockStorage) Close() error {
	return nil
}

func TestRateLimiter_IP(t *testing.T) {
	// Configuração de teste
	cfg := &config.Config{
		IPRateLimit:    5,
		TokenRateLimit: 20,
		BlockDuration:  300,
	}

	mockStorage := NewMockStorage()
	service := limiter.NewService(mockStorage, cfg)

	ctx := context.Background()
	ip := "192.168.1.1"

	// Testa que as primeiras 5 requisições são permitidas
	for i := 0; i < 5; i++ {
		info, err := service.Allow(ctx, ip, "")
		if err != nil {
			t.Fatalf("Erro não esperado: %v", err)
		}
		if !info.Allowed {
			t.Fatalf("Requisição %d deveria ser permitida", i+1)
		}
	}

	// A sexta requisição deve ser bloqueada
	info, err := service.Allow(ctx, ip, "")
	if err != nil {
		t.Fatalf("Erro não esperado: %v", err)
	}
	if info.Allowed {
		t.Fatal("A sexta requisição deveria ser bloqueada")
	}
}

func TestRateLimiter_Token(t *testing.T) {
	// Configuração de teste
	cfg := &config.Config{
		IPRateLimit:    5,
		TokenRateLimit: 10,
		BlockDuration:  300,
	}

	mockStorage := NewMockStorage()
	service := limiter.NewService(mockStorage, cfg)

	ctx := context.Background()
	ip := "192.168.1.1"
	token := "abc123"

	// Testa que as primeiras 10 requisições com token são permitidas
	for i := 0; i < 10; i++ {
		info, err := service.Allow(ctx, ip, token)
		if err != nil {
			t.Fatalf("Erro não esperado: %v", err)
		}
		if !info.Allowed {
			t.Fatalf("Requisição %d com token deveria ser permitida", i+1)
		}
	}

	// A 11ª requisição com token deve ser bloqueada
	info, err := service.Allow(ctx, ip, token)
	if err != nil {
		t.Fatalf("Erro não esperado: %v", err)
	}
	if info.Allowed {
		t.Fatal("A 11ª requisição com token deveria ser bloqueada")
	}
}

// Add this new test to verify that the mock is working correctly
func TestMockStorage_VerifyBlocking(t *testing.T) {
	// Create storage and verify its behavior directly
	mockStorage := NewMockStorage()

	ctx := context.Background()
	key := "test-key"

	// Initial state should not be blocked
	blocked, err := mockStorage.IsBlocked(ctx, key)
	if err != nil {
		t.Fatalf("Error checking blocked status: %v", err)
	}
	if blocked {
		t.Fatal("Key should not be blocked initially")
	}

	// Block the key and verify
	err = mockStorage.Block(ctx, key, time.Second)
	if err != nil {
		t.Fatalf("Error blocking key: %v", err)
	}

	blocked, err = mockStorage.IsBlocked(ctx, key)
	if err != nil {
		t.Fatalf("Error checking blocked status: %v", err)
	}
	if !blocked {
		t.Fatal("Key should be blocked after Block() call")
	}

	// Reset counter before testing increment
	mockStorage.counters[key] = 0

	// Verify counter increments
	count, err := mockStorage.Increment(ctx, key)
	if err != nil {
		t.Fatalf("Error incrementing counter: %v", err)
	}
	if count != 1 {
		t.Fatalf("Expected count to be 1, got %d", count)
	}
}

// This test uses a different configuration to avoid caching issues
func TestRateLimiter_TokenWithDifferentLimit(t *testing.T) {
	// Use a different token limit to force a new test run
	cfg := &config.Config{
		IPRateLimit:    5,
		TokenRateLimit: 3, // Much smaller limit to make it obvious
		BlockDuration:  300,
	}

	mockStorage := NewMockStorage()
	service := limiter.NewService(mockStorage, cfg)

	ctx := context.Background()
	ip := "192.168.1.1"
	token := "different-token"

	t.Logf("Token rate limit set to: %d", cfg.TokenRateLimit)

	// Make exactly the configured number of allowed requests
	for i := 0; i < cfg.TokenRateLimit; i++ {
		info, err := service.Allow(ctx, ip, token)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		t.Logf("Request %d: Allowed=%v, Count=%d/%d",
			i+1, info.Allowed, info.CurrentCount, info.Limit)
		if !info.Allowed {
			t.Fatalf("Request %d should be allowed (count: %d, limit: %d)",
				i+1, info.CurrentCount, info.Limit)
		}
	}

	// The next request should be blocked
	info, err := service.Allow(ctx, ip, token)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	t.Logf("Final request: Allowed=%v, Count=%d/%d",
		info.Allowed, info.CurrentCount, info.Limit)

	if info.Allowed {
		t.Fatalf("Request should be blocked (count: %d, limit: %d)",
			info.CurrentCount, info.Limit)
	}
}
