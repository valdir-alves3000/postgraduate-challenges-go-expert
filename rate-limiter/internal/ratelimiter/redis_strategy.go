package ratelimiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(addr string) *RedisStore {
	client := redis.NewClient(&redis.Options{Addr: addr})
	return &RedisStore{client: client}
}

func (r *RedisStore) Get(key string) (int, error) {
	val, err := r.client.Get(context.Background(), key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

func (r *RedisStore) Increment(key string, expiration int) (int, error) {
	pipe := r.client.TxPipeline()
	incr := pipe.Incr(context.Background(), key)
	pipe.Expire(context.Background(), key, time.Duration(expiration)*time.Second)
	_, err := pipe.Exec(context.Background())

	return int(incr.Val()), err
}

func (r *RedisStore) Block(key string, duration int) error {
	return r.client.Set(context.Background(), key+":blocked", 1, time.Duration(duration)*time.Second).Err()
}

func (r *RedisStore) IsBlocked(key string) (bool, error) {
	val, err := r.client.Get(context.Background(), key+":blocked").Result()
	if err == redis.Nil {
		return false, nil
	}
	if val == "1" {
		return true, nil
	}

	return false, err
}

func (r *RedisStore) ListKeys(pattern string) ([]string, error) {
	keys, err := r.client.Keys(context.Background(), pattern).Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}
