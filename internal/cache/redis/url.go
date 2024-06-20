package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"url-storage/internal/cache"
	"url-storage/internal/config"
	"url-storage/internal/storage"
)

type Storage interface {
	Insert(ctx context.Context, url string, alias string) (int64, error)
	Get(ctx context.Context, alias string) (string, error)
}

type Cache struct {
	client  *redis.Client
	storage Storage
}

func New(cfg *config.RedisConfig, s Storage) (*Cache, error) {
	const op = "redis.New"

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Cache{
		client:  client,
		storage: s,
	}, nil
}

func (c *Cache) SetURL(ctx context.Context, alias string, url string) (int64, error) {
	const op = "redis.Set"

	id, err := c.storage.Insert(ctx, alias, url)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExist) {
			return -1, cache.ErrAlreadyExist
		}
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	err = c.client.Set(ctx, alias, url, 0).Err()
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (c *Cache) GetURL(ctx context.Context, alias string) (string, error) {
	const op = "redis.Get"

	url, err := c.client.Get(ctx, alias).Result()
	if err != nil {
		if err == redis.Nil {
			url, err = c.storage.Get(ctx, alias)
			if err != nil {
				if errors.Is(err, storage.ErrNotFound) {
					return "", fmt.Errorf("%s: %w", op, cache.ErrNotFound)
				}
				return "", fmt.Errorf("%s: %w", op, err)
			}

			err = c.client.Set(ctx, alias, url, 0).Err()
			if err != nil {
				return "", fmt.Errorf("%s: %w", op, err)
			}
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}
