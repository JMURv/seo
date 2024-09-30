package redis

import (
	"context"
	"encoding/json"
	"fmt"
	errs "github.com/JMURv/par-pro-seo/internal/cache"
	controller "github.com/JMURv/par-pro-seo/internal/controller"
	cfg "github.com/JMURv/par-pro-seo/pkg/config"
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"log"
	"time"
)

type Cache struct {
	cli *redis.Client
}

func New(conf *cfg.RedisConfig) controller.CacheRepo {
	cli := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Pass,
		DB:       0,
	})
	_, err := cli.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return &Cache{cli: cli}
}

func (c *Cache) Close() {
	if err := c.cli.Close(); err != nil {
		zap.L().Debug("Failed to close connection to Redis: ", zap.Error(err))
	}
}

func (c *Cache) Get(ctx context.Context, key string, dest any) error {
	const op = "Get"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	val, err := c.cli.Get(ctx, key).Bytes()
	if err == redis.Nil {
		zap.L().Debug("[CACHE] MISS", zap.String("key", key))
		return errs.ErrNotFoundInCache
	} else if err != nil {
		zap.L().Debug("[CACHE] ERROR", zap.String("key", key), zap.Error(err))
		return err
	}

	if err = json.Unmarshal(val, dest); err != nil {
		zap.L().Debug("[CACHE] ERROR", zap.String("key", key), zap.Error(err))
		return err
	}

	zap.L().Debug("[CACHE] HIT", zap.String("key", key))
	return nil
}

func (c *Cache) Set(ctx context.Context, t time.Duration, key string, val any) error {
	const op = "Set"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	if err := c.cli.Set(ctx, key, val, t).Err(); err != nil {
		zap.L().Debug("[CACHE] ERROR", zap.String("key", key), zap.Error(err))
		return err
	}

	zap.L().Debug("[CACHE] SET", zap.String("key", key))
	return nil
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	const op = "Delete"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	if err := c.cli.Del(ctx, key).Err(); err != nil {
		zap.L().Debug("[CACHE] ERROR", zap.String("key", key), zap.Error(err))
		return err
	}
	return nil
}

func (c *Cache) InvalidateKeysByPattern(ctx context.Context, pattern string) error {
	var cursor uint64
	for {
		var keys []string
		var err error

		keys, cursor, err = c.cli.Scan(ctx, cursor, pattern, 100).Result() // 100 keys at a time
		if err != nil {
			return fmt.Errorf("failed to scan redis: %w", err)
		}

		if len(keys) > 0 {
			if err := c.cli.Del(ctx, keys...).Err(); err != nil {
				return fmt.Errorf("failed to delete keys: %w", err)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}
