package limiter

import (
	"context"
	"github.com/joho/godotenv"
	"os"
	"testing"
)

// Função init será executada antes dos testes
func init() {
	// Carregar configurações de teste
	if err := godotenv.Load("../.env.test"); err != nil {
		// Fall back to default values if .env.test is not found
		os.Setenv("IP_MAX_REQUESTS", "5")
		os.Setenv("IP_BLOCK_DURATION", "5")
		os.Setenv("TOKEN_MAX_REQUESTS", "10")
		os.Setenv("TOKEN_BLOCK_DURATION", "5")
	}
}

type MockStorage struct {
	counts  map[string]int
	blocked map[string]bool
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		counts:  make(map[string]int),
		blocked: make(map[string]bool),
	}
}

func (m *MockStorage) Increment(key string) (int, error) {
	m.counts[key]++
	return m.counts[key], nil
}

func (m *MockStorage) Reset(key string) error {
	delete(m.counts, key)
	delete(m.blocked, key)
	return nil
}

func (m *MockStorage) IsBlocked(key string) (bool, error) {
	return m.blocked[key], nil
}

func (m *MockStorage) Block(key string, duration int) error {
	m.blocked[key] = true
	return nil
}

func TestRateLimiter(t *testing.T) {
	storage := NewMockStorage()
	limiter := NewRateLimiter(storage)

	t.Run("IP Limiting", func(t *testing.T) {
		storage.Reset("ip:192.168.1.1")
		ip := "192.168.1.1"

		// Primeiras 5 requisições devem ser permitidas
		for i := 0; i < 5; i++ {
			allowed, err := limiter.IsAllowed(ip, "")
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !allowed {
				t.Errorf("Request %d should be allowed", i+1)
			}
		}

		// Sexta requisição deve ser bloqueada
		allowed, err := limiter.IsAllowed(ip, "")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if allowed {
			t.Error("Sixth request should be blocked")
		}
	})

	t.Run("Token Limiting", func(t *testing.T) {
		storage.Reset("token:abc123")
		token := "abc123"

		// Primeiras 10 requisições devem ser permitidas
		for i := 0; i < 10; i++ {
			allowed, err := limiter.IsAllowed("", token)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !allowed {
				t.Errorf("Request %d should be allowed", i+1)
			}
		}

		// 11ª requisição deve ser bloqueada
		allowed, err := limiter.IsAllowed("", token)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if allowed {
			t.Error("11th request should be blocked")
		}
	})
}

func TestRedisStorage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Redis integration test in short mode")
	}

	storage := NewRedisStorage("redis://localhost:6379/0")

	// Tentar um ping no Redis
	ctx := context.Background()
	if err := storage.client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping integration tests")
	}

	t.Run("Increment and Block", func(t *testing.T) {
		key := "test-key"

		// Limpar chave antes do teste
		storage.Reset(key)

		count, err := storage.Increment(key)
		if err != nil {
			t.Fatalf("Failed to increment: %v", err)
		}
		if count != 1 {
			t.Errorf("Expected count 1, got %d", count)
		}

		err = storage.Block(key, 1)
		if err != nil {
			t.Fatalf("Failed to block: %v", err)
		}

		blocked, err := storage.IsBlocked(key)
		if err != nil {
			t.Fatalf("Failed to check blocked status: %v", err)
		}
		if !blocked {
			t.Error("Expected key to be blocked")
		}
	})
}
