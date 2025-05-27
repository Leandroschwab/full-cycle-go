package storage

import (
	"context"
	"time"
)

// Storage define a interface para os mecanismos de armazenamento
type Storage interface {
	// Incrementa o contador para uma chave e retorna o valor atual
	Increment(ctx context.Context, key string) (int, error)

	// Verifica se uma chave está bloqueada
	IsBlocked(ctx context.Context, key string) (bool, error)

	// Bloqueia uma chave por uma duração específica
	Block(ctx context.Context, key string, duration time.Duration) error

	// Reseta o contador para uma chave
	Reset(ctx context.Context, key string) error

	// Fecha a conexão com o armazenamento
	Close() error
}
