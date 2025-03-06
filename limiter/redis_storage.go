package limiter

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/sony/gobreaker"
	"time"
)

type RedisStorage struct {
	client  *redis.Client
	ctx     context.Context
	breaker *gobreaker.CircuitBreaker
}

func NewRedisStorage(url string) *RedisStorage {
	opts, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opts)

	settings := gobreaker.Settings{
		Name:        "redis",
		MaxRequests: 3,
		Interval:    10 * time.Second,
		Timeout:     60 * time.Second,
	}

	return &RedisStorage{
		client:  client,
		ctx:     context.Background(),
		breaker: gobreaker.NewCircuitBreaker(settings),
	}
}

func (r *RedisStorage) Increment(key string) (int, error) {
	result, err := r.breaker.Execute(func() (interface{}, error) {
		pipe := r.client.Pipeline()
		incr := pipe.Incr(r.ctx, key)
		pipe.Expire(r.ctx, key, time.Second)
		_, err := pipe.Exec(r.ctx)
		if err != nil {
			return 0, err
		}
		return incr.Val(), nil
	})

	if err != nil {
		return 0, err
	}

	return int(result.(int64)), nil
}

func (r *RedisStorage) Reset(key string) error {
	_, err := r.breaker.Execute(func() (interface{}, error) {
		return nil, r.client.Del(r.ctx, key).Err()
	})
	return err
}

func (r *RedisStorage) IsBlocked(key string) (bool, error) {
	result, err := r.breaker.Execute(func() (interface{}, error) {
		return r.client.Exists(r.ctx, "blocked:"+key).Result()
	})

	if err != nil {
		return false, err
	}

	return result.(int64) == 1, nil
}

func (r *RedisStorage) Block(key string, duration int) error {
	_, err := r.breaker.Execute(func() (interface{}, error) {
		return nil, r.client.Set(r.ctx, "blocked:"+key, 1, time.Duration(duration)*time.Minute).Err()
	})
	return err
}

func (r *RedisStorage) Close() error {
	return r.client.Close()
}
