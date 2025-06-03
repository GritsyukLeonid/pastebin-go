package logging

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLogger struct {
	client *redis.Client
	ttl    time.Duration
}

type Logger interface {
	LogChange(entity, id, action string) error
}

func NewRedisLogger(addr string, ttl time.Duration) *RedisLogger {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisLogger{
		client: client,
		ttl:    ttl,
	}
}

func (r *RedisLogger) LogChange(entity, id, action string) error {
	ctx := context.Background()
	key := fmt.Sprintf("log:%s:%s:%d", entity, id, time.Now().Unix())
	value := action
	return r.client.Set(ctx, key, value, r.ttl).Err()
}
