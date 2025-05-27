package storage

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	counterPrefix = "rate_limit:counter:"
	blockedPrefix = "rate_limit:blocked:"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(redisURL string) *RedisStorage {
	client := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	return &RedisStorage{
		client: client,
	}
}

func (s *RedisStorage) Increment(ctx context.Context, key string) (int, error) {
	counterKey := counterPrefix + key

	// Verifica se a chave já existe
	exists, err := s.client.Exists(ctx, counterKey).Result()
	if err != nil {
		return 0, err
	}

	// Se a chave não existe, cria com TTL de 1 segundo
	if exists == 0 {
		err = s.client.Set(ctx, counterKey, 0, time.Second).Err()
		if err != nil {
			return 0, err
		}
	}

	// Incrementa o contador
	count, err := s.client.Incr(ctx, counterKey).Result()
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (s *RedisStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	blockedKey := blockedPrefix + key

	// Verifica se a chave está na lista de bloqueados
	exists, err := s.client.Exists(ctx, blockedKey).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}

func (s *RedisStorage) Block(ctx context.Context, key string, duration time.Duration) error {
	blockedKey := blockedPrefix + key

	// Adiciona a chave à lista de bloqueados com o TTL especificado
	return s.client.Set(ctx, blockedKey, "1", duration).Err()
}

func (s *RedisStorage) Reset(ctx context.Context, key string) error {
	counterKey := counterPrefix + key

	// Remove o contador
	return s.client.Del(ctx, counterKey).Err()
}

func (s *RedisStorage) Close() error {
	return s.client.Close()
}
