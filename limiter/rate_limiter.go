package limiter

import (
	"fmt"
	"os"
	"strconv"
)

type RateLimiterInterface interface {
	IsAllowed(ip string, token string) (bool, error)
}

type RateLimiter struct {
	storage Storage
}

func NewRateLimiter(storage Storage) RateLimiterInterface {
	return &RateLimiter{
		storage: storage,
	}
}

func (r *RateLimiter) IsAllowed(ip string, token string) (bool, error) {
	// Verificar se está bloqueado
	if token != "" {
		blocked, err := r.storage.IsBlocked(fmt.Sprintf("token:%s", token))
		if err != nil {
			return false, err
		}
		if blocked {
			return false, nil
		}
		return r.checkToken(token)
	}

	blocked, err := r.storage.IsBlocked(fmt.Sprintf("ip:%s", ip))
	if err != nil {
		return false, err
	}
	if blocked {
		return false, nil
	}
	return r.checkIP(ip)
}

func (r *RateLimiter) checkIP(ip string) (bool, error) {
	key := fmt.Sprintf("ip:%s", ip)
	count, err := r.storage.Increment(key)
	if err != nil {
		return false, err
	}

	maxRequests, _ := strconv.Atoi(os.Getenv("IP_MAX_REQUESTS"))
	if maxRequests == 0 {
		maxRequests = 5 // valor padrão
	}

	if count > maxRequests {
		blockDuration, _ := strconv.Atoi(os.Getenv("IP_BLOCK_DURATION"))
		if blockDuration == 0 {
			blockDuration = 5 // valor padrão
		}
		err = r.storage.Block(key, blockDuration)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

func (r *RateLimiter) checkToken(token string) (bool, error) {
	key := fmt.Sprintf("token:%s", token)
	count, err := r.storage.Increment(key)
	if err != nil {
		return false, err
	}

	maxRequests, _ := strconv.Atoi(os.Getenv("TOKEN_MAX_REQUESTS"))
	if maxRequests == 0 {
		maxRequests = 10 // valor padrão
	}

	if count > maxRequests {
		blockDuration, _ := strconv.Atoi(os.Getenv("TOKEN_BLOCK_DURATION"))
		if blockDuration == 0 {
			blockDuration = 5 // valor padrão
		}
		err = r.storage.Block(key, blockDuration)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}
